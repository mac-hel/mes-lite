package csvimport

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
)

// ProductionEntryRecord is a validated and typed historical production-entry import row.
type ProductionEntryRecord struct {
	RowNumber   int
	EmployeeID  string
	ProductSKU  string
	Quantity    int
	Workstation string
	Timestamp   time.Time
	Comment     string
}

// RowError describes one validation error found in one CSV row.
type RowError struct {
	RowNumber int
	Field     string
	Message   string
}

// Error returns a human-readable validation error message.
func (e RowError) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("row %d: %s", e.RowNumber, e.Message)
	}
	return fmt.Sprintf("row %d field %s: %s", e.RowNumber, e.Field, e.Message)
}

// ValidationResult contains valid rows and row-level validation errors from a CSV stream.
type ValidationResult struct {
	Records []ProductionEntryRecord
	Errors  []RowError
}

// ValidateProductionEntries reads the whole CSV stream row by row and collects row validation errors.
func ValidateProductionEntries(reader *ProductionEntryReader) (ValidationResult, error) {
	var result ValidationResult

	for {
		raw, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result, nil
			}
			return ValidationResult{}, err
		}

		record, rowErrors := validateProductionEntryRow(raw)
		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, rowErrors...)
			continue
		}
		result.Records = append(result.Records, record)
	}
}

func validateProductionEntryRow(row ProductionEntryRow) (ProductionEntryRecord, []RowError) {
	var rowErrors []RowError

	if row.EmployeeID == "" {
		rowErrors = append(rowErrors, newRowError(row.RowNumber, "employee_id", "is required"))
	}
	if row.ProductSKU == "" {
		rowErrors = append(rowErrors, newRowError(row.RowNumber, "product_sku", "is required"))
	}
	if row.Workstation == "" {
		rowErrors = append(rowErrors, newRowError(row.RowNumber, "workstation", "is required"))
	}

	quantity, ok := parseQuantity(row, &rowErrors)
	timestamp, timestampOK := parseTimestamp(row, &rowErrors)
	if len(rowErrors) > 0 {
		return ProductionEntryRecord{}, rowErrors
	}
	if !ok || !timestampOK {
		return ProductionEntryRecord{}, rowErrors
	}

	return ProductionEntryRecord{
		RowNumber:   row.RowNumber,
		EmployeeID:  row.EmployeeID,
		ProductSKU:  row.ProductSKU,
		Quantity:    quantity,
		Workstation: row.Workstation,
		Timestamp:   timestamp.UTC(),
		Comment:     row.Comment,
	}, nil
}

func parseQuantity(row ProductionEntryRow, rowErrors *[]RowError) (int, bool) {
	if row.Quantity == "" {
		*rowErrors = append(*rowErrors, newRowError(row.RowNumber, "quantity", "is required"))
		return 0, false
	}

	quantity, err := strconv.Atoi(row.Quantity)
	if err != nil {
		*rowErrors = append(*rowErrors, newRowError(row.RowNumber, "quantity", "must be an integer"))
		return 0, false
	}
	if quantity <= 0 {
		*rowErrors = append(*rowErrors, newRowError(row.RowNumber, "quantity", "must be greater than zero"))
		return 0, false
	}
	if quantity > math.MaxInt32 {
		*rowErrors = append(*rowErrors, newRowError(row.RowNumber, "quantity", "must fit PostgreSQL integer"))
		return 0, false
	}

	return quantity, true
}

func parseTimestamp(row ProductionEntryRow, rowErrors *[]RowError) (time.Time, bool) {
	if row.Timestamp == "" {
		*rowErrors = append(*rowErrors, newRowError(row.RowNumber, "timestamp", "is required"))
		return time.Time{}, false
	}

	timestamp, err := time.Parse(time.RFC3339, row.Timestamp)
	if err != nil {
		*rowErrors = append(*rowErrors, newRowError(row.RowNumber, "timestamp", "must be RFC3339"))
		return time.Time{}, false
	}

	return timestamp, true
}

func newRowError(rowNumber int, field, message string) RowError {
	return RowError{
		RowNumber: rowNumber,
		Field:     field,
		Message:   message,
	}
}
