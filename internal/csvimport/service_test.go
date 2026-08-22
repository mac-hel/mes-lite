package csvimport

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/mac-hel/mes-lite/internal/production"
)

func TestServiceImportProductionEntries_ImportsValidRowsAndReportsValidationErrors(t *testing.T) {
	store := NewInMemoryStore()
	service := NewService(store)
	csv := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
		",sku-1,0,,not-time,invalid",
		"emp-2,sku-2,7,ws-2,2026-08-20T11:00:00Z,also valid",
	}, "\n")

	summary, err := service.ImportProductionEntries(t.Context(), strings.NewReader(csv))
	if err != nil {
		t.Fatalf("ImportProductionEntries() error = %v", err)
	}

	if summary.TotalRows != 3 || summary.ValidRows != 2 || summary.InvalidRows != 1 || summary.ImportedRows != 2 {
		t.Fatalf("summary = %+v", summary)
	}
	if len(summary.Errors) != 4 {
		t.Fatalf("errors length = %d, want 4: %+v", len(summary.Errors), summary.Errors)
	}
	if len(store.Entries()) != 2 {
		t.Fatalf("stored entries length = %d, want 2", len(store.Entries()))
	}
}

func TestServiceImportProductionEntries_ReturnsFatalCSVErrors(t *testing.T) {
	service := NewService(NewInMemoryStore())

	_, err := service.ImportProductionEntries(t.Context(), strings.NewReader("employee_id,product_sku\n"))
	if !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("ImportProductionEntries() error = %v, want ErrInvalidHeader", err)
	}
}

func TestServiceImportProductionEntries_ReportsBatchErrorsWithoutImportingRows(t *testing.T) {
	store := &rowFailureStore{failedRows: map[int]error{3: production.ErrInvalidEntry}}
	service := NewService(store)
	csv := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
		"missing,sku-1,7,ws-2,2026-08-20T11:00:00Z,bad reference",
	}, "\n")

	summary, err := service.ImportProductionEntries(t.Context(), strings.NewReader(csv))
	if err != nil {
		t.Fatalf("ImportProductionEntries() error = %v", err)
	}

	if summary.TotalRows != 2 || summary.ValidRows != 1 || summary.InvalidRows != 1 || summary.ImportedRows != 1 {
		t.Fatalf("summary = %+v", summary)
	}
	if len(summary.Errors) != 1 {
		t.Fatalf("errors length = %d, want 1: %+v", len(summary.Errors), summary.Errors)
	}
	if summary.Errors[0].RowNumber != 3 || summary.Errors[0].Message != production.ErrInvalidEntry.Error() {
		t.Fatalf("error = %+v", summary.Errors[0])
	}
	if len(store.savedRows) != 1 || store.savedRows[0] != 2 {
		t.Fatalf("saved rows = %v, want [2]", store.savedRows)
	}
}

func TestServiceImportProductionEntries_ContinuesAfterAlreadyExistsOnRetry(t *testing.T) {
	store := &rowFailureStore{failedRows: map[int]error{2: production.ErrAlreadyExists}}
	service := NewService(store)
	csv := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,already imported",
		"emp-2,sku-2,7,ws-2,2026-08-20T11:00:00Z,new row",
	}, "\n")

	summary, err := service.ImportProductionEntries(t.Context(), strings.NewReader(csv))
	if err != nil {
		t.Fatalf("ImportProductionEntries() error = %v", err)
	}

	if summary.TotalRows != 2 || summary.ValidRows != 1 || summary.InvalidRows != 1 || summary.ImportedRows != 1 {
		t.Fatalf("summary = %+v", summary)
	}
	if len(summary.Errors) != 1 {
		t.Fatalf("errors length = %d, want 1: %+v", len(summary.Errors), summary.Errors)
	}
	if summary.Errors[0].RowNumber != 2 || !strings.Contains(summary.Errors[0].Message, production.ErrAlreadyExists.Error()) {
		t.Fatalf("error = %+v", summary.Errors[0])
	}
	if len(store.savedRows) != 1 || store.savedRows[0] != 3 {
		t.Fatalf("saved rows = %v, want [3]", store.savedRows)
	}
}

func TestServiceImportProductionEntries_ReturnsUnexpectedStoreErrors(t *testing.T) {
	service := NewService(unexpectedErrorStore{})
	csv := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
	}, "\n")

	_, err := service.ImportProductionEntries(t.Context(), strings.NewReader(csv))
	if !errors.Is(err, errUnexpectedStore) {
		t.Fatalf("ImportProductionEntries() error = %v, want errUnexpectedStore", err)
	}
}

func TestServiceImportProductionEntries_SavesLargeInputInBatches(t *testing.T) {
	store := &recordingStore{}
	service := NewService(store)
	var input strings.Builder
	input.WriteString("employee_id,product_sku,quantity,workstation,timestamp,comment\n")
	for i := range 1205 {
		_, _ = fmt.Fprintf(&input, "emp-1,sku-1,%d,ws-1,2026-08-20T10:00:00Z,row %d\n", i+1, i+1)
	}

	summary, err := service.ImportProductionEntries(t.Context(), strings.NewReader(input.String()))
	if err != nil {
		t.Fatalf("ImportProductionEntries() error = %v", err)
	}

	if summary.TotalRows != 1205 || summary.ValidRows != 1205 || summary.InvalidRows != 0 || summary.ImportedRows != 1205 {
		t.Fatalf("summary = %+v", summary)
	}
	wantBatchSizes := []int{500, 500, 205}
	if len(store.batchSizes) != len(wantBatchSizes) {
		t.Fatalf("batch sizes = %v, want %v", store.batchSizes, wantBatchSizes)
	}
	for i := range wantBatchSizes {
		if store.batchSizes[i] != wantBatchSizes[i] {
			t.Fatalf("batch sizes = %v, want %v", store.batchSizes, wantBatchSizes)
		}
	}
}

func TestServiceImportProductionEntries_CapsReportedErrorsForLargeInvalidInput(t *testing.T) {
	store := &recordingStore{}
	service := NewService(store)
	var input strings.Builder
	input.WriteString("employee_id,product_sku,quantity,workstation,timestamp,comment\n")
	for range 1200 {
		input.WriteString(",sku-1,0,ws-1,2026-08-20T10:00:00Z,invalid\n")
	}

	summary, err := service.ImportProductionEntries(t.Context(), strings.NewReader(input.String()))
	if err != nil {
		t.Fatalf("ImportProductionEntries() error = %v", err)
	}

	if summary.TotalRows != 1200 || summary.ValidRows != 0 || summary.InvalidRows != 1200 || summary.ImportedRows != 0 {
		t.Fatalf("summary = %+v", summary)
	}
	if len(summary.Errors) != maxReportedErrors {
		t.Fatalf("reported errors length = %d, want %d", len(summary.Errors), maxReportedErrors)
	}
	if !summary.ErrorLimitReached {
		t.Fatal("ErrorLimitReached = false, want true")
	}
	if len(store.batchSizes) != 0 {
		t.Fatalf("batch sizes = %v, want no save calls", store.batchSizes)
	}
}

var errUnexpectedStore = errors.New("unexpected store failure")

type unexpectedErrorStore struct{}

func (unexpectedErrorStore) SaveBatch(context.Context, []ProductionEntryRecord) ([]production.Entry, error) {
	return nil, errUnexpectedStore
}

type recordingStore struct {
	batchSizes []int
}

func (s *recordingStore) SaveBatch(_ context.Context, records []ProductionEntryRecord) ([]production.Entry, error) {
	s.batchSizes = append(s.batchSizes, len(records))
	entries := make([]production.Entry, len(records))
	return entries, nil
}

type rowFailureStore struct {
	failedRows map[int]error
	savedRows  []int
}

func (s *rowFailureStore) SaveBatch(_ context.Context, records []ProductionEntryRecord) ([]production.Entry, error) {
	for _, record := range records {
		if err, ok := s.failedRows[record.RowNumber]; ok {
			return nil, BatchError{RowNumber: record.RowNumber, Err: err}
		}
	}

	for _, record := range records {
		s.savedRows = append(s.savedRows, record.RowNumber)
	}
	entries := make([]production.Entry, len(records))
	return entries, nil
}
