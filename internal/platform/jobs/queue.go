package jobs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// DefaultQueueCapacity bounds how many jobs can wait for a worker.
const DefaultQueueCapacity = 128

// ErrQueueFull is returned when the queue has no room for another waiting job.
var ErrQueueFull = errors.New("background job queue is full")

// ErrQueueClosed is returned when the queue no longer accepts or hands out work.
var ErrQueueClosed = errors.New("background job queue is closed")

// ErrNotFound is returned when a background job cannot be found by ID.
var ErrNotFound = errors.New("background job not found")

// ErrAlreadyExists is returned when a job ID is already tracked by the queue.
var ErrAlreadyExists = errors.New("background job already exists")

// ErrInvalidStatusTransition is returned when a job cannot move between two lifecycle states.
var ErrInvalidStatusTransition = errors.New("invalid background job status transition")

// ErrInvalidProgress is returned when reported job progress is outside 0..100.
var ErrInvalidProgress = errors.New("invalid background job progress")

// Queue is an in-memory FIFO queue of background jobs.
//
// It owns two pieces of state with different jobs to do:
//
//   - jobs is the canonical status of every tracked job, guarded by mu.
//   - pending is the handoff channel between producers and consumers. It carries
//     job IDs rather than job values, so a consumer always reads current state
//     from the map instead of acting on a stale copy taken at enqueue time.
//
// The queue never closes pending. There are many producers, and closing a
// channel that other goroutines may still send on is a panic, not a shutdown
// signal. Shutdown uses a dedicated done channel instead, which has exactly one
// closer: Close.
type Queue struct {
	mu      sync.RWMutex
	jobs    map[string]Job
	pending chan string
	done    chan struct{}
	closed  bool
}

// NewQueue creates an in-memory queue that can hold capacity waiting jobs.
// A capacity of zero or less falls back to [DefaultQueueCapacity].
func NewQueue(capacity int) *Queue {
	if capacity <= 0 {
		capacity = DefaultQueueCapacity
	}

	return &Queue{
		jobs:    make(map[string]Job, capacity),
		pending: make(chan string, capacity),
		done:    make(chan struct{}),
	}
}

// Enqueue adds a queued job and makes it available to consumers.
//
// It never blocks. A full queue returns [ErrQueueFull] so an HTTP caller gets
// backpressure immediately instead of waiting on a channel that may stay full
// for minutes. The buffer capacity is the burst the application is willing to
// absorb.
func (q *Queue) Enqueue(ctx context.Context, job Job) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := job.Validate(); err != nil {
		return err
	}
	if job.Status != StatusQueued {
		return fmt.Errorf("job %q must be queued, got %q: %w", job.ID, job.Status, ErrInvalidJob)
	}

	// The lock is held across the channel send on purpose. The send is
	// non-blocking, so it cannot deadlock, and holding the lock means a consumer
	// that receives the ID cannot look it up before the map entry exists.
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return fmt.Errorf("enqueue job %q: %w", job.ID, ErrQueueClosed)
	}
	if _, ok := q.jobs[job.ID]; ok {
		return fmt.Errorf("background job %q: %w", job.ID, ErrAlreadyExists)
	}

	select {
	case q.pending <- job.ID:
		q.jobs[job.ID] = job.clone()
		return nil
	default:
		return fmt.Errorf("capacity %d reached: %w", cap(q.pending), ErrQueueFull)
	}
}

// Dequeue returns the next waiting job, blocking until one is available, the
// context is cancelled or the queue is closed and drained.
//
// Waiting jobs are drained before the close signal is honoured so a shutdown
// does not silently discard work that was already accepted.
func (q *Queue) Dequeue(ctx context.Context) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}

	// Fast path: waiting work always wins over the close signal.
	if id, ok := q.tryReceive(); ok {
		return q.Find(ctx, id)
	}

	select {
	case id := <-q.pending:
		return q.Find(ctx, id)
	case <-q.done:
		// A select with several ready cases picks one at random, so reaching this
		// case does not prove the queue is empty. Enqueue can no longer send once
		// done is closed, which makes one final look exact.
		if id, ok := q.tryReceive(); ok {
			return q.Find(ctx, id)
		}
		return Job{}, ErrQueueClosed
	case <-ctx.Done():
		return Job{}, ctx.Err()
	}
}

// tryReceive takes one waiting job ID without blocking.
func (q *Queue) tryReceive() (string, bool) {
	select {
	case id := <-q.pending:
		return id, true
	default:
		return "", false
	}
}

// Find returns a copy of one tracked job. Returns [ErrNotFound] if unknown.
func (q *Queue) Find(_ context.Context, id string) (Job, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}

	return job.clone(), nil
}

// MarkRunning moves a queued job to running and records its start time.
func (q *Queue) MarkRunning(ctx context.Context, id string, startedAt time.Time) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}
	if job.Status != StatusQueued {
		return Job{}, fmt.Errorf("job %q cannot move from %q to %q: %w", id, job.Status, StatusRunning, ErrInvalidStatusTransition)
	}

	job.Status = StatusRunning
	job.StartedAt = startedAt.UTC()
	job.FinishedAt = time.Time{}
	job.Error = ""
	q.jobs[id] = job

	return job.clone(), nil
}

// MarkSucceeded moves a running job to succeeded and records its finish time.
func (q *Queue) MarkSucceeded(ctx context.Context, id string, finishedAt time.Time) (Job, error) {
	return q.markFinished(ctx, id, StatusSucceeded, "", finishedAt)
}

// MarkCancelled moves a queued or running job to cancelled and records its finish time.
func (q *Queue) MarkCancelled(ctx context.Context, id string, finishedAt time.Time) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}
	if job.Status.Terminal() {
		return job.clone(), nil
	}

	job.Status = StatusCancelled
	job.CancelRequested = true
	job.FinishedAt = finishedAt.UTC()
	q.jobs[id] = job

	return job.clone(), nil
}

// MarkFailed moves a running job to failed and records the failure message.
func (q *Queue) MarkFailed(ctx context.Context, id string, message string, finishedAt time.Time) (Job, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "job failed"
	}
	return q.markFinished(ctx, id, StatusFailed, message, finishedAt)
}

func (q *Queue) markFinished(ctx context.Context, id string, status Status, message string, finishedAt time.Time) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}
	if job.Status != StatusRunning {
		return Job{}, fmt.Errorf("job %q cannot move from %q to %q: %w", id, job.Status, status, ErrInvalidStatusTransition)
	}

	job.Status = status
	if status == StatusSucceeded {
		job.Progress = 100
	}
	job.FinishedAt = finishedAt.UTC()
	job.Error = message
	q.jobs[id] = job

	return job.clone(), nil
}

// ReportProgress records the latest progress percentage for a running job.
func (q *Queue) ReportProgress(ctx context.Context, id string, progress int) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}
	if progress < 0 || progress > 100 {
		return Job{}, fmt.Errorf("progress %d must be between 0 and 100: %w", progress, ErrInvalidProgress)
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}
	if job.Status != StatusRunning {
		return Job{}, fmt.Errorf("job %q cannot report progress while %q: %w", id, job.Status, ErrInvalidStatusTransition)
	}

	job.Progress = progress
	q.jobs[id] = job

	return job.clone(), nil
}

// RecordResult stores handler-produced result data for a running job.
func (q *Queue) RecordResult(ctx context.Context, id string, result []byte) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}
	if job.Status != StatusRunning {
		return Job{}, fmt.Errorf("job %q cannot record result while %q: %w", id, job.Status, ErrInvalidStatusTransition)
	}

	job.Result = copyPayload(result)
	q.jobs[id] = job

	return job.clone(), nil
}

// RequestCancellation records a cancellation request. Queued jobs are cancelled
// immediately; running jobs must observe their context and stop cooperatively.
func (q *Queue) RequestCancellation(ctx context.Context, id string, requestedAt time.Time) (Job, error) {
	if err := ctx.Err(); err != nil {
		return Job{}, err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("background job %q: %w", id, ErrNotFound)
	}
	if job.Status.Terminal() {
		return job.clone(), nil
	}

	job.CancelRequested = true
	if job.Status == StatusQueued {
		job.Status = StatusCancelled
		job.FinishedAt = requestedAt.UTC()
	}
	q.jobs[id] = job

	return job.clone(), nil
}

// CancellationRequested reports whether a job has an active cancellation request.
func (q *Queue) CancellationRequested(ctx context.Context, id string) (bool, error) {
	job, err := q.Find(ctx, id)
	if err != nil {
		return false, err
	}
	return job.CancelRequested, nil
}

// Len returns how many jobs are waiting for a consumer.
func (q *Queue) Len() int {
	return len(q.pending)
}

// Capacity returns how many jobs can wait for a consumer at once.
func (q *Queue) Capacity() int {
	return cap(q.pending)
}

// Close stops the queue from accepting new work and releases consumers that are
// waiting for jobs. It is safe to call more than once.
func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}
	q.closed = true
	close(q.done)
}
