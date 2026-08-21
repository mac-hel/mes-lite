package csvimport

import (
	"context"
	"errors"
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
	store := NewFailingInMemoryStore(BatchError{RowNumber: 3, Err: production.ErrInvalidEntry})
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

	if summary.TotalRows != 2 || summary.ValidRows != 1 || summary.InvalidRows != 1 || summary.ImportedRows != 0 {
		t.Fatalf("summary = %+v", summary)
	}
	if len(summary.Errors) != 1 {
		t.Fatalf("errors length = %d, want 1: %+v", len(summary.Errors), summary.Errors)
	}
	if summary.Errors[0].RowNumber != 3 || summary.Errors[0].Message != production.ErrInvalidEntry.Error() {
		t.Fatalf("error = %+v", summary.Errors[0])
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

var errUnexpectedStore = errors.New("unexpected store failure")

type unexpectedErrorStore struct{}

func (unexpectedErrorStore) SaveBatch(context.Context, []ProductionEntryRecord) ([]production.Entry, error) {
	return nil, errUnexpectedStore
}
