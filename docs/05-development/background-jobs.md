# Background Jobs Development Notes

This document shows how to use the in-memory job queue and worker pool.

The important rule is separation of responsibility:

- `jobs.Queue` owns job state.
- `jobs.WorkerPool` owns worker goroutines and cancellation wiring.
- job handlers own business work.
- cancellation is cooperative: handlers must check `ctx.Done()` or `ctx.Err()`.

## Enqueue A Job

```go
queue := jobs.NewQueue(jobs.DefaultQueueCapacity)

job, err := jobs.NewJob(
	ids.New(),
	jobs.TypeProductionEntryImport,
	[]byte("csv payload or future payload reference"),
	time.Now(),
)
if err != nil {
	return err
}

if err := queue.Enqueue(ctx, job); err != nil {
	return err
}
```

`Enqueue` is non-blocking. If the queue is full, it returns `jobs.ErrQueueFull` so the caller can apply backpressure instead of waiting indefinitely.

## Start A Worker Pool

```go
handlers := map[jobs.Type]jobs.Handler{
	jobs.TypeProductionEntryImport: importProductionEntries,
}

workers, err := jobs.NewWorkerPool(queue, 4, handlers)
if err != nil {
	return err
}

workers.Start(ctx)
defer func() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = workers.Stop(shutdownCtx)
}()
```

The worker count is the maximum number of jobs that can run at the same time. Queue capacity is different: it controls how many jobs may wait before workers pick them up.

## Write A Handler With Progress

```go
func importProductionEntries(ctx context.Context, job jobs.Job) error {
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := processOneBatch(ctx, job.Payload); err != nil {
			return err
		}

		_, err := queue.ReportProgress(ctx, job.ID, i+1)
		if err != nil {
			return err
		}
	}

	return nil
}
```

The handler receives a `jobs.Job` copy. It must not mutate job state directly. Report state changes through queue methods such as `ReportProgress`.

## Cancel A Job

```go
job, err := workers.Cancel(ctx, jobID)
if err != nil {
	return err
}

fmt.Println(job.Status, job.CancelRequested)
```

Cancelling a queued job marks it `cancelled` immediately.

Cancelling a running job records `CancelRequested` and calls the job's context cancel function. That closes `ctx.Done()` for the handler. The goroutine is not killed by force.

## Cleanup On Cancellation

Use `defer` for cleanup, then return `ctx.Err()` when cancellation is observed.

```go
func importWithCleanup(ctx context.Context, job jobs.Job) error {
	file, err := os.Open("production.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row, err := readNextRow(file)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		if err := insertRow(ctx, tx, row); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
```

If `workers.Cancel` is called while this handler runs:

- the job context is cancelled
- `<-ctx.Done()` becomes ready
- the handler returns `context.Canceled`
- deferred cleanup runs
- the worker records the final job status as `cancelled`

## Bad Handler Example

```go
func badHandler(ctx context.Context, job jobs.Job) error {
	time.Sleep(10 * time.Minute)
	return nil
}
```

This handler ignores cancellation. Calling `workers.Cancel` records the request, but the job will not stop until `Sleep` finishes.

Prefer cancellable operations:

```go
func betterHandler(ctx context.Context, job jobs.Job) error {
	select {
	case <-time.After(10 * time.Minute):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

## Common Mistakes

- Do not close the queue's internal work channel from a handler.
- Do not mutate a `jobs.Job` copy and expect the queue to see it.
- Do not launch one goroutine per job unless the workload has its own separate concurrency limit.
- Do not ignore `ctx.Done()` in long-running handlers.
- Do not treat cancellation as failure when it was requested by the user.
