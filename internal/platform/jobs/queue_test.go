package jobs

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func newTestJob(t *testing.T, id string) Job {
	t.Helper()

	job, err := NewJob(id, TypeProductionEntryImport, []byte("payload"), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return job
}

func TestQueue_EnqueueDequeueIsFIFO(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(4)

	for _, id := range []string{"job-1", "job-2", "job-3"} {
		if err := queue.Enqueue(ctx, newTestJob(t, id)); err != nil {
			t.Fatalf("enqueue %s: %v", id, err)
		}
	}

	if queue.Len() != 3 {
		t.Errorf("expected 3 waiting jobs, got %d", queue.Len())
	}

	for _, want := range []string{"job-1", "job-2", "job-3"} {
		job, err := queue.Dequeue(ctx)
		if err != nil {
			t.Fatalf("dequeue: %v", err)
		}
		if job.ID != want {
			t.Errorf("expected %q, got %q", want, job.ID)
		}
		if job.Status != StatusQueued {
			t.Errorf("expected dequeued job to still be queued, got %q", job.Status)
		}
	}

	if queue.Len() != 0 {
		t.Errorf("expected an empty queue, got %d waiting jobs", queue.Len())
	}
}

func TestQueue_EnqueueRejectsDuplicateID(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(4)

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	err := queue.Enqueue(ctx, newTestJob(t, "job-1"))
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
	if queue.Len() != 1 {
		t.Errorf("expected the duplicate to be rejected, got %d waiting jobs", queue.Len())
	}
}

func TestQueue_EnqueueRejectsNonQueuedJob(t *testing.T) {
	job := newTestJob(t, "job-1")
	job.Status = StatusRunning

	err := NewQueue(1).Enqueue(t.Context(), job)
	if !errors.Is(err, ErrInvalidJob) {
		t.Fatalf("expected ErrInvalidJob, got %v", err)
	}
}

func TestQueue_EnqueueReturnsQueueFullWithoutBlocking(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	err := queue.Enqueue(ctx, newTestJob(t, "job-2"))
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}

	// A rejected job must not stay tracked, otherwise its ID could never be reused.
	if _, err := queue.Find(ctx, "job-2"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected the rejected job to be untracked, got %v", err)
	}
}

func TestQueue_DequeueBlocksUntilWorkArrives(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)

	type result struct {
		job Job
		err error
	}
	results := make(chan result, 1)

	go func() {
		job, err := queue.Dequeue(ctx)
		results <- result{job: job, err: err}
	}()

	select {
	case got := <-results:
		t.Fatalf("expected Dequeue to block on an empty queue, got %+v", got)
	case <-time.After(20 * time.Millisecond):
	}

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-results:
		if got.err != nil {
			t.Fatalf("dequeue: %v", got.err)
		}
		if got.job.ID != "job-1" {
			t.Errorf("expected job-1, got %q", got.job.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("expected Dequeue to return once a job was enqueued")
	}
}

func TestQueue_DequeueReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	queue := NewQueue(1)

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err := queue.Dequeue(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestQueue_CloseDrainsWaitingJobsBeforeReportingClosed(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(2)

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}
	queue.Close()

	job, err := queue.Dequeue(ctx)
	if err != nil {
		t.Fatalf("expected the accepted job to be drained, got %v", err)
	}
	if job.ID != "job-1" {
		t.Errorf("expected job-1, got %q", job.ID)
	}

	if _, err := queue.Dequeue(ctx); !errors.Is(err, ErrQueueClosed) {
		t.Fatalf("expected ErrQueueClosed, got %v", err)
	}
}

func TestQueue_CloseRejectsNewWorkAndIsIdempotent(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(2)

	queue.Close()
	queue.Close()

	err := queue.Enqueue(ctx, newTestJob(t, "job-1"))
	if !errors.Is(err, ErrQueueClosed) {
		t.Fatalf("expected ErrQueueClosed, got %v", err)
	}
}

func TestQueue_CloseReleasesWaitingConsumers(t *testing.T) {
	queue := NewQueue(1)
	errs := make(chan error, 1)

	go func() {
		_, err := queue.Dequeue(t.Context())
		errs <- err
	}()

	time.Sleep(20 * time.Millisecond)
	queue.Close()

	select {
	case err := <-errs:
		if !errors.Is(err, ErrQueueClosed) {
			t.Fatalf("expected ErrQueueClosed, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("expected Close to release the waiting consumer")
	}
}

func TestQueue_FindReturnsACopy(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	job, err := queue.Find(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	job.Status = StatusFailed
	job.Payload[0] = 'X'

	tracked, err := queue.Find(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != StatusQueued {
		t.Errorf("expected tracked status to stay queued, got %q", tracked.Status)
	}
	if string(tracked.Payload) != "payload" {
		t.Errorf("expected tracked payload to stay unchanged, got %q", tracked.Payload)
	}
}

func TestQueue_FindUnknownJob(t *testing.T) {
	if _, err := NewQueue(1).Find(t.Context(), "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// TestQueue_ConcurrentProducersAndConsumers is written for `go test -race`:
// every enqueued job must be handed to exactly one consumer.
func TestQueue_ConcurrentProducersAndConsumers(t *testing.T) {
	ctx := t.Context()

	const (
		producers        = 8
		jobsPerProducer  = 25
		consumers        = 4
		totalJobs        = producers * jobsPerProducer
		consumerCapacity = totalJobs
	)

	queue := NewQueue(consumerCapacity)

	var producerGroup sync.WaitGroup
	for p := range producers {
		producerGroup.Add(1)
		go func() {
			defer producerGroup.Done()
			for i := range jobsPerProducer {
				// t.Fatal must not be called from a non-test goroutine, so this
				// builds the job inline instead of using the helper.
				job, err := NewJob(fmt.Sprintf("job-%d-%d", p, i), TypeProductionEntryImport, []byte("payload"), time.Now())
				if err != nil {
					t.Errorf("new job: %v", err)
					return
				}
				if err := queue.Enqueue(ctx, job); err != nil {
					t.Errorf("enqueue: %v", err)
					return
				}
			}
		}()
	}

	dequeued := make(chan string, totalJobs)

	var consumerGroup sync.WaitGroup
	for range consumers {
		consumerGroup.Add(1)
		go func() {
			defer consumerGroup.Done()
			for {
				job, err := queue.Dequeue(ctx)
				if err != nil {
					if !errors.Is(err, ErrQueueClosed) {
						t.Errorf("dequeue: %v", err)
					}
					return
				}
				dequeued <- job.ID
			}
		}()
	}

	producerGroup.Wait()
	queue.Close()
	consumerGroup.Wait()
	close(dequeued)

	seen := make(map[string]int, totalJobs)
	for id := range dequeued {
		seen[id]++
	}

	if len(seen) != totalJobs {
		t.Fatalf("expected %d distinct jobs, got %d", totalJobs, len(seen))
	}
	for id, count := range seen {
		if count != 1 {
			t.Errorf("job %q was delivered %d times, want exactly once", id, count)
		}
	}
}
