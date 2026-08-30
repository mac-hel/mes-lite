package jobs

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewWorkerPool_ValidatesConfiguration(t *testing.T) {
	if _, err := NewWorkerPool(nil, 1, nil); !errors.Is(err, ErrInvalidWorkerPool) {
		t.Fatalf("expected missing queue to return ErrInvalidWorkerPool, got %v", err)
	}

	_, err := NewWorkerPool(NewQueue(1), 1, map[Type]Handler{
		TypeProductionEntryImport: nil,
	})
	if !errors.Is(err, ErrInvalidWorkerPool) {
		t.Fatalf("expected nil handler to return ErrInvalidWorkerPool, got %v", err)
	}
}

func TestWorkerPool_ExecutesJobSuccessfully(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)
	executed := make(chan Job, 1)

	pool, err := NewWorkerPool(queue, 1, map[Type]Handler{
		TypeProductionEntryImport: func(_ context.Context, job Job) error {
			executed <- job
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}
	pool.Start(ctx)
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}

	select {
	case job := <-executed:
		if job.ID != "job-1" {
			t.Errorf("expected job-1, got %q", job.ID)
		}
		if job.Status != StatusRunning {
			t.Errorf("expected handler to receive a running job, got %q", job.Status)
		}
	case <-time.After(time.Second):
		t.Fatal("expected handler to execute")
	}

	tracked, err := queue.Find(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != StatusSucceeded {
		t.Fatalf("expected succeeded job, got %q", tracked.Status)
	}
	if tracked.StartedAt.IsZero() || tracked.FinishedAt.IsZero() {
		t.Fatalf("expected lifecycle timestamps to be set, got started=%s finished=%s", tracked.StartedAt, tracked.FinishedAt)
	}
}

func TestWorkerPool_RecordsHandlerFailure(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)

	pool, err := NewWorkerPool(queue, 1, map[Type]Handler{
		TypeProductionEntryImport: func(context.Context, Job) error {
			return errors.New("import failed")
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}
	pool.Start(ctx)
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}

	tracked, err := queue.Find(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != StatusFailed {
		t.Fatalf("expected failed job, got %q", tracked.Status)
	}
	if tracked.Error != "import failed" {
		t.Fatalf("expected failure message to be recorded, got %q", tracked.Error)
	}
}

func TestWorkerPool_UnknownJobTypeFailsJob(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)
	pool, err := NewWorkerPool(queue, 1, nil)
	if err != nil {
		t.Fatal(err)
	}

	job, err := NewJob("job-1", Type("unknown"), nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(ctx, job); err != nil {
		t.Fatal(err)
	}

	pool.Start(ctx)
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}

	tracked, err := queue.Find(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != StatusFailed {
		t.Fatalf("expected failed job, got %q", tracked.Status)
	}
	if tracked.Error == "" {
		t.Fatal("expected unknown job type to be recorded as an error")
	}
}

func TestWorkerPool_MultipleWorkersExecuteEachJobOnce(t *testing.T) {
	ctx := t.Context()

	const totalJobs = 100
	queue := NewQueue(totalJobs)

	var mu sync.Mutex
	executed := make(map[string]int, totalJobs)
	pool, err := NewWorkerPool(queue, 4, map[Type]Handler{
		TypeProductionEntryImport: func(_ context.Context, job Job) error {
			mu.Lock()
			defer mu.Unlock()
			executed[job.ID]++
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := range totalJobs {
		if err := queue.Enqueue(ctx, newTestJob(t, fmt.Sprintf("job-%03d", i))); err != nil {
			t.Fatal(err)
		}
	}

	pool.Start(ctx)
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}

	if len(executed) != totalJobs {
		t.Fatalf("expected %d executed jobs, got %d", totalJobs, len(executed))
	}
	for id, count := range executed {
		if count != 1 {
			t.Errorf("job %q executed %d times, want once", id, count)
		}
	}
}

func TestWorkerPool_StopWaitsForRunningWork(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)
	started := make(chan struct{})
	release := make(chan struct{})

	pool, err := NewWorkerPool(queue, 1, map[Type]Handler{
		TypeProductionEntryImport: func(context.Context, Job) error {
			close(started)
			<-release
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	pool.Start(ctx)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("expected job to start")
	}

	stopped := make(chan error, 1)
	go func() {
		stopped <- stopPool(t, pool)
	}()

	select {
	case err := <-stopped:
		t.Fatalf("expected Stop to wait for running work, got %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	close(release)

	select {
	case err := <-stopped:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("expected Stop to return after running work finished")
	}
}

func TestWorkerPool_StopIsIdempotent(t *testing.T) {
	pool, err := NewWorkerPool(NewQueue(1), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	pool.Start(t.Context())

	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}
}

func TestWorkerPool_CancelQueuedJob(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)
	pool, err := NewWorkerPool(queue, 1, map[Type]Handler{
		TypeProductionEntryImport: func(context.Context, Job) error {
			t.Fatal("cancelled queued job should not execute")
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	job, err := pool.Cancel(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != StatusCancelled {
		t.Fatalf("expected queued job to be cancelled, got %q", job.Status)
	}

	pool.Start(ctx)
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}
}

func TestWorkerPool_CancelRunningJob(t *testing.T) {
	ctx := t.Context()
	queue := NewQueue(1)
	started := make(chan struct{})

	pool, err := NewWorkerPool(queue, 1, map[Type]Handler{
		TypeProductionEntryImport: func(ctx context.Context, job Job) error {
			if _, err := queue.ReportProgress(ctx, job.ID, 50); err != nil {
				return err
			}
			close(started)
			<-ctx.Done()
			return ctx.Err()
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(ctx, newTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	pool.Start(ctx)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("expected job to start")
	}

	job, err := pool.Cancel(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != StatusRunning || !job.CancelRequested {
		t.Fatalf("expected running job with cancellation requested, got %+v", job)
	}
	if err := stopPool(t, pool); err != nil {
		t.Fatal(err)
	}

	tracked, err := queue.Find(ctx, "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != StatusCancelled {
		t.Fatalf("expected cancelled job, got %q", tracked.Status)
	}
	if tracked.Progress != 50 {
		t.Fatalf("expected progress to be preserved, got %d", tracked.Progress)
	}
}

func stopPool(t *testing.T, pool *WorkerPool) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	return pool.Stop(ctx)
}
