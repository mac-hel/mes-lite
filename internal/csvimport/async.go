package csvimport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mac-hel/mes-lite/internal/platform/ids"
	"github.com/mac-hel/mes-lite/internal/platform/jobs"
)

// JobQueue is the queue behavior needed by the async CSV import workflow.
type JobQueue interface {
	Enqueue(context.Context, jobs.Job) error
}

// ProgressRecorder records job progress and result data from the worker handler.
type ProgressRecorder interface {
	ReportProgress(context.Context, string, int) (jobs.Job, error)
	RecordResult(context.Context, string, []byte) (jobs.Job, error)
}

// AsyncService accepts import uploads and turns them into background jobs.
type AsyncService struct {
	queue   JobQueue
	tempDir string
	now     func() time.Time
}

type productionEntryImportPayload struct {
	Path string `json:"path"`
}

// NewAsyncService creates an async CSV import service.
func NewAsyncService(queue JobQueue, tempDir string) *AsyncService {
	return &AsyncService{queue: queue, tempDir: tempDir, now: time.Now}
}

// EnqueueProductionEntries copies a CSV upload to a temporary file and enqueues an import job.
func (s *AsyncService) EnqueueProductionEntries(ctx context.Context, input io.Reader) (jobs.Job, error) {
	if err := ctx.Err(); err != nil {
		return jobs.Job{}, err
	}

	file, err := os.CreateTemp(s.tempDir, "mes-lite-production-import-*.csv")
	if err != nil {
		return jobs.Job{}, fmt.Errorf("create import temp file: %w", err)
	}
	path := file.Name()
	removeTemp := true
	defer func() {
		_ = file.Close()
		if removeTemp {
			_ = os.Remove(path)
		}
	}()

	if _, err := io.Copy(file, input); err != nil {
		return jobs.Job{}, fmt.Errorf("write import temp file: %w", err)
	}
	if err := file.Close(); err != nil {
		return jobs.Job{}, fmt.Errorf("close import temp file: %w", err)
	}

	payload, err := json.Marshal(productionEntryImportPayload{Path: path})
	if err != nil {
		return jobs.Job{}, fmt.Errorf("encode import job payload: %w", err)
	}
	job, err := jobs.NewJob(ids.New(), jobs.TypeProductionEntryImport, payload, s.now())
	if err != nil {
		return jobs.Job{}, err
	}
	if err := s.queue.Enqueue(ctx, job); err != nil {
		return jobs.Job{}, err
	}

	removeTemp = false
	return job, nil
}

// NewProductionEntriesJobHandler creates a background handler for queued production-entry imports.
func NewProductionEntriesJobHandler(service *Service, recorder ProgressRecorder) jobs.Handler {
	return func(ctx context.Context, job jobs.Job) error {
		var payload productionEntryImportPayload
		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return fmt.Errorf("decode import job payload: %w", err)
		}
		if payload.Path == "" {
			return fmt.Errorf("import job payload path is required")
		}

		file, err := os.Open(payload.Path)
		if err != nil {
			return fmt.Errorf("open import temp file: %w", err)
		}
		defer func() {
			_ = file.Close()
			_ = os.Remove(payload.Path)
		}()

		if _, err := recorder.ReportProgress(ctx, job.ID, 5); err != nil {
			return err
		}
		summary, err := service.ImportProductionEntries(ctx, file)
		if err != nil {
			return err
		}

		result, err := json.Marshal(summary)
		if err != nil {
			return fmt.Errorf("encode import summary: %w", err)
		}
		if _, err := recorder.RecordResult(ctx, job.ID, result); err != nil {
			return err
		}
		_, err = recorder.ReportProgress(ctx, job.ID, 95)
		return err
	}
}
