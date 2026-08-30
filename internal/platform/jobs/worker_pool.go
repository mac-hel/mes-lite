package jobs

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// DefaultWorkerCount is used when a worker pool is created with a non-positive worker count.
const DefaultWorkerCount = 1

// ErrInvalidWorkerPool is returned when worker pool configuration is unusable.
var ErrInvalidWorkerPool = errors.New("invalid background worker pool")

// Handler executes one background job.
type Handler func(context.Context, Job) error

// WorkerPool runs queued jobs with a fixed number of worker goroutines.
type WorkerPool struct {
	queue       *Queue
	workerCount int
	handlers    map[Type]Handler
	now         func() time.Time

	mu        sync.Mutex
	running   map[string]context.CancelFunc
	startOnce sync.Once
	stopOnce  sync.Once
	wg        sync.WaitGroup
	cancel    context.CancelFunc
}

// NewWorkerPool creates a worker pool that executes jobs from queue.
func NewWorkerPool(queue *Queue, workerCount int, handlers map[Type]Handler) (*WorkerPool, error) {
	if queue == nil {
		return nil, fmt.Errorf("queue is required: %w", ErrInvalidWorkerPool)
	}
	if workerCount <= 0 {
		workerCount = DefaultWorkerCount
	}

	copiedHandlers := make(map[Type]Handler, len(handlers))
	for jobType, handler := range handlers {
		if handler == nil {
			return nil, fmt.Errorf("handler for job type %q is nil: %w", jobType, ErrInvalidWorkerPool)
		}
		copiedHandlers[jobType] = handler
	}

	return &WorkerPool{
		queue:       queue,
		workerCount: workerCount,
		handlers:    copiedHandlers,
		now:         time.Now,
		running:     make(map[string]context.CancelFunc),
	}, nil
}

// Start launches worker goroutines. Calling Start more than once has no effect.
func (p *WorkerPool) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	p.startOnce.Do(func() {
		workerCtx, cancel := context.WithCancel(ctx)
		p.cancel = cancel

		for range p.workerCount {
			p.wg.Add(1)
			go p.work(workerCtx)
		}
	})
}

// Stop closes the queue, waits for workers to finish accepted work and returns
// ctx.Err if the caller stops waiting before the workers exit.
func (p *WorkerPool) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	p.stopOnce.Do(func() {
		p.queue.Close()
	})

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		if p.cancel != nil {
			p.cancel()
		}
		return ctx.Err()
	}
}

// Cancel requests cancellation for one job. Running jobs receive context
// cancellation; queued jobs are marked cancelled before a worker can start them.
func (p *WorkerPool) Cancel(ctx context.Context, id string) (Job, error) {
	job, err := p.queue.RequestCancellation(ctx, id, p.now())
	if err != nil {
		return Job{}, err
	}
	if job.Status == StatusRunning {
		p.mu.Lock()
		cancel := p.running[id]
		p.mu.Unlock()
		if cancel != nil {
			// This does not kill the worker goroutine. It closes the job context's
			// Done channel so a cooperative handler can notice cancellation, clean up
			// with defers and return context.Canceled.
			cancel()
		}
	}
	return job, nil
}

func (p *WorkerPool) work(ctx context.Context) {
	defer p.wg.Done()

	for {
		job, err := p.queue.Dequeue(ctx)
		if err != nil {
			return
		}

		p.execute(ctx, job.ID)
	}
}

func (p *WorkerPool) execute(ctx context.Context, id string) {
	job, err := p.queue.MarkRunning(ctx, id, p.now())
	if err != nil {
		return
	}
	// Each job gets its own child context. The worker context stops the whole
	// pool; this child context lets Cancel stop one running job without stopping
	// every worker.
	jobCtx, cancel := context.WithCancel(ctx)
	p.mu.Lock()
	// Store only the cancel function, not the job itself. The queue remains the
	// owner of job state; this map is just the cancellation wiring for currently
	// running handlers.
	p.running[job.ID] = cancel
	p.mu.Unlock()
	defer func() {
		// Always call cancel when the job finishes, even after success. It releases
		// context resources and wakes any child operation still waiting on Done.
		cancel()
		p.mu.Lock()
		// After this point a user can no longer cancel this job through the running
		// map. Terminal state is already recorded in the queue.
		delete(p.running, job.ID)
		p.mu.Unlock()
	}()

	handler := p.handlers[job.Type]
	if handler == nil {
		_, _ = p.queue.MarkFailed(context.Background(), job.ID, fmt.Sprintf("no handler registered for job type %q", job.Type), p.now())
		return
	}

	if err := handler(jobCtx, job.clone()); err != nil {
		if errors.Is(err, context.Canceled) && p.cancelRequested(job.ID) {
			// A handler that returns context.Canceled after a recorded cancel request
			// completed cooperatively, so the final job state is cancelled, not failed.
			_, _ = p.queue.MarkCancelled(context.Background(), job.ID, p.now())
			return
		}
		_, _ = p.queue.MarkFailed(context.Background(), job.ID, err.Error(), p.now())
		return
	}
	if p.cancelRequested(job.ID) {
		// The handler may return nil after noticing cancellation and doing cleanup.
		// The queue's cancellation flag still decides the final state.
		_, _ = p.queue.MarkCancelled(context.Background(), job.ID, p.now())
		return
	}

	_, _ = p.queue.MarkSucceeded(context.Background(), job.ID, p.now())
}

func (p *WorkerPool) cancelRequested(id string) bool {
	requested, err := p.queue.CancellationRequested(context.Background(), id)
	return err == nil && requested
}
