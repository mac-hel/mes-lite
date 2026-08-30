package csvimport

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mac-hel/mes-lite/internal/platform/jobs"
)

func TestAsyncServiceEnqueueProductionEntries_EnqueuesJobWithTempFilePayload(t *testing.T) {
	queue := jobs.NewQueue(1)
	service := NewAsyncService(queue, t.TempDir())
	csv := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
	}, "\n")

	job, err := service.EnqueueProductionEntries(t.Context(), strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if job.Type != jobs.TypeProductionEntryImport || job.Status != jobs.StatusQueued {
		t.Fatalf("unexpected job: %+v", job)
	}

	var payload productionEntryImportPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	stored, err := os.ReadFile(payload.Path)
	if err != nil {
		t.Fatal(err)
	}
	if string(stored) != csv {
		t.Fatalf("stored upload = %q, want %q", stored, csv)
	}

	queued, err := queue.Find(t.Context(), job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if queued.ID != job.ID {
		t.Fatalf("queued job ID = %q, want %q", queued.ID, job.ID)
	}
}

func TestAsyncServiceEnqueueProductionEntries_RemovesTempFileWhenQueueIsFull(t *testing.T) {
	queue := jobs.NewQueue(1)
	if err := queue.Enqueue(t.Context(), newCSVImportJob(t, "existing", nil)); err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	service := NewAsyncService(queue, tempDir)

	_, err := service.EnqueueProductionEntries(t.Context(), strings.NewReader("employee_id,product_sku,quantity,workstation,timestamp,comment\n"))
	if !errors.Is(err, jobs.ErrQueueFull) {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected temp file cleanup, got %d files", len(entries))
	}
}

func TestProductionEntriesJobHandler_ImportsCSVRecordsResultAndRemovesTempFile(t *testing.T) {
	queue := jobs.NewQueue(1)
	store := NewInMemoryStore()
	service := NewService(store)
	path := writeTempCSV(t, strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
	}, "\n"))
	payload, err := json.Marshal(productionEntryImportPayload{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	job := newCSVImportJob(t, "job-1", payload)
	if err := queue.Enqueue(t.Context(), job); err != nil {
		t.Fatal(err)
	}

	workers, err := jobs.NewWorkerPool(queue, 1, map[jobs.Type]jobs.Handler{
		jobs.TypeProductionEntryImport: NewProductionEntriesJobHandler(service, queue),
	})
	if err != nil {
		t.Fatal(err)
	}
	workers.Start(t.Context())
	if err := stopCSVImportWorkers(t, workers); err != nil {
		t.Fatal(err)
	}

	tracked, err := queue.Find(t.Context(), "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != jobs.StatusSucceeded {
		t.Fatalf("expected succeeded job, got %+v", tracked)
	}
	if tracked.Progress != 100 {
		t.Fatalf("expected progress 100, got %d", tracked.Progress)
	}
	var summary ImportSummary
	if err := json.Unmarshal(tracked.Result, &summary); err != nil {
		t.Fatal(err)
	}
	if summary.ImportedRows != 1 || len(store.Entries()) != 1 {
		t.Fatalf("summary = %+v, entries = %d", summary, len(store.Entries()))
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected temp file to be removed, stat error = %v", err)
	}
}

func TestProductionEntriesJobHandler_FailsJobForInvalidCSV(t *testing.T) {
	queue := jobs.NewQueue(1)
	service := NewService(NewInMemoryStore())
	path := writeTempCSV(t, "employee_id,quantity\n")
	payload, err := json.Marshal(productionEntryImportPayload{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(t.Context(), newCSVImportJob(t, "job-1", payload)); err != nil {
		t.Fatal(err)
	}

	workers, err := jobs.NewWorkerPool(queue, 1, map[jobs.Type]jobs.Handler{
		jobs.TypeProductionEntryImport: NewProductionEntriesJobHandler(service, queue),
	})
	if err != nil {
		t.Fatal(err)
	}
	workers.Start(t.Context())
	if err := stopCSVImportWorkers(t, workers); err != nil {
		t.Fatal(err)
	}

	tracked, err := queue.Find(t.Context(), "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if tracked.Status != jobs.StatusFailed {
		t.Fatalf("expected failed job, got %+v", tracked)
	}
	if !strings.Contains(tracked.Error, ErrInvalidHeader.Error()) {
		t.Fatalf("expected invalid header error, got %q", tracked.Error)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected temp file to be removed, stat error = %v", err)
	}
}

func newCSVImportJob(t *testing.T, id string, payload []byte) jobs.Job {
	t.Helper()

	job, err := jobs.NewJob(id, jobs.TypeProductionEntryImport, payload, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return job
}

func writeTempCSV(t *testing.T, content string) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "import-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	path := file.Name()
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func stopCSVImportWorkers(t *testing.T, workers *jobs.WorkerPool) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	return workers.Stop(ctx)
}
