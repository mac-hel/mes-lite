package csvimport

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

func TestValidateProductionEntries_ReturnsTypedRecords(t *testing.T) {
	reader := newTestProductionEntryReader(t, strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T12:30:00+02:00,done",
	}, "\n"))

	result, err := ValidateProductionEntries(reader)
	if err != nil {
		t.Fatalf("ValidateProductionEntries() error = %v", err)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("errors = %+v, want none", result.Errors)
	}
	if len(result.Records) != 1 {
		t.Fatalf("records length = %d, want 1", len(result.Records))
	}

	record := result.Records[0]
	if record.RowNumber != 2 || record.EmployeeID != "emp-1" || record.ProductSKU != "sku-1" || record.Quantity != 12 || record.Workstation != "ws-1" || record.Comment != "done" {
		t.Fatalf("record = %+v", record)
	}
	wantTimestamp := time.Date(2026, 8, 20, 10, 30, 0, 0, time.UTC)
	if !record.Timestamp.Equal(wantTimestamp) || record.Timestamp.Location() != time.UTC {
		t.Fatalf("timestamp = %v, want %v in UTC", record.Timestamp, wantTimestamp)
	}
}

func TestValidateProductionEntries_CollectsRowErrorsAndKeepsValidRows(t *testing.T) {
	reader := newTestProductionEntryReader(t, strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
		",,0,,not-time,invalid",
		"emp-2,sku-2,7,ws-2,2026-08-20T11:00:00Z,also valid",
	}, "\n"))

	result, err := ValidateProductionEntries(reader)
	if err != nil {
		t.Fatalf("ValidateProductionEntries() error = %v", err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records length = %d, want 2", len(result.Records))
	}
	if len(result.Errors) != 5 {
		t.Fatalf("errors length = %d, want 5: %+v", len(result.Errors), result.Errors)
	}

	want := []RowError{
		{RowNumber: 3, Field: "employee_id", Message: "is required"},
		{RowNumber: 3, Field: "product_sku", Message: "is required"},
		{RowNumber: 3, Field: "workstation", Message: "is required"},
		{RowNumber: 3, Field: "quantity", Message: "must be greater than zero"},
		{RowNumber: 3, Field: "timestamp", Message: "must be RFC3339"},
	}
	for i := range want {
		if result.Errors[i] != want[i] {
			t.Fatalf("error[%d] = %+v, want %+v", i, result.Errors[i], want[i])
		}
	}
}

func TestValidateProductionEntries_RejectsInvalidQuantityValues(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		message string
	}{
		{name: "missing", value: "", message: "is required"},
		{name: "not integer", value: "12.5", message: "must be an integer"},
		{name: "negative", value: "-1", message: "must be greater than zero"},
		{name: "too large", value: "2147483648", message: "must fit PostgreSQL integer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := newTestProductionEntryReader(t, strings.Join([]string{
				"employee_id,product_sku,quantity,workstation,timestamp,comment",
				"emp-1,sku-1," + tt.value + ",ws-1,2026-08-20T10:00:00Z,comment",
			}, "\n"))

			result, err := ValidateProductionEntries(reader)
			if err != nil {
				t.Fatalf("ValidateProductionEntries() error = %v", err)
			}
			if len(result.Records) != 0 {
				t.Fatalf("records length = %d, want 0", len(result.Records))
			}
			if len(result.Errors) != 1 {
				t.Fatalf("errors length = %d, want 1: %+v", len(result.Errors), result.Errors)
			}
			if result.Errors[0] != (RowError{RowNumber: 2, Field: "quantity", Message: tt.message}) {
				t.Fatalf("error = %+v", result.Errors[0])
			}
		})
	}
}

func TestValidateProductionEntries_ReturnsFatalReadError(t *testing.T) {
	reader := newTestProductionEntryReader(t, strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
		"emp-2,sku-2,7,ws-2,2026-08-20T11:00:00Z",
	}, "\n"))

	_, err := ValidateProductionEntries(reader)
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("ValidateProductionEntries() error = %v, want ErrInvalidRecord", err)
	}
}

func TestRowError_Error(t *testing.T) {
	err := RowError{RowNumber: 12, Field: "quantity", Message: "must be an integer"}
	if err.Error() != "row 12 field quantity: must be an integer" {
		t.Fatalf("Error() = %q", err.Error())
	}

	err = RowError{RowNumber: 12, Message: "could not parse"}
	if err.Error() != "row 12: could not parse" {
		t.Fatalf("Error() = %q", err.Error())
	}
}

func newTestProductionEntryReader(t *testing.T, input string) *ProductionEntryReader {
	t.Helper()

	reader, err := NewProductionEntryReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("NewProductionEntryReader() error = %v", err)
	}
	return reader
}

func TestValidateProductionEntries_EmptyDataReturnsEmptySlices(t *testing.T) {
	reader := newTestProductionEntryReader(t, "employee_id,product_sku,quantity,workstation,timestamp,comment\n")

	result, err := ValidateProductionEntries(reader)
	if err != nil {
		t.Fatalf("ValidateProductionEntries() error = %v", err)
	}
	if len(result.Records) != 0 {
		t.Fatalf("records length = %d, want 0", len(result.Records))
	}
	if len(result.Errors) != 0 {
		t.Fatalf("errors length = %d, want 0", len(result.Errors))
	}
}

func TestProductionEntryReader_InvalidFieldCountMapsToImportError(t *testing.T) {
	reader := newTestProductionEntryReader(t, strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z",
	}, "\n"))

	_, err := reader.Read()
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("Read() error = %v, want ErrInvalidRecord", err)
	}
	if errors.Is(err, io.EOF) {
		t.Fatalf("Read() error = %v, did not want io.EOF", err)
	}
}
