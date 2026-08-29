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

	handler := p.handlers[job.Type]
	if handler == nil {
		_, _ = p.queue.MarkFailed(context.Background(), job.ID, fmt.Sprintf("no handler registered for job type %q", job.Type), p.now())
		return
	}

	if err := handler(ctx, job.clone()); err != nil {
		_, _ = p.queue.MarkFailed(context.Background(), job.ID, err.Error(), p.now())
		return
	}

	_, _ = p.queue.MarkSucceeded(context.Background(), job.ID, p.now())
}
