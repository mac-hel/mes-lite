package csvimport

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var (
	// ErrInvalidHeader is returned when the CSV header does not match the import contract.
	ErrInvalidHeader = errors.New("invalid csv header")
	// ErrInvalidRecord is returned when a CSV record has the wrong shape.
	ErrInvalidRecord = errors.New("invalid csv record")
)

var productionEntryHeader = []string{
	"employee_id",
	"product_sku",
	"quantity",
	"workstation",
	"timestamp",
	"comment",
}

// ProductionEntryRow is one raw CSV data row for a historical production entry.
// Field parsing and business validation are intentionally handled by later import steps.
type ProductionEntryRow struct {
	RowNumber   int
	EmployeeID  string
	ProductSKU  string
	Quantity    string
	Workstation string
	Timestamp   string
	Comment     string
}

// ProductionEntryReader streams historical production-entry rows from CSV input.
type ProductionEntryReader struct {
	reader    *csv.Reader
	rowNumber int
}

// NewProductionEntryReader creates a streaming reader for the production-entry CSV import format.
func NewProductionEntryReader(r io.Reader) (*ProductionEntryReader, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("missing csv header: %w", ErrInvalidHeader)
		}
		return nil, fmt.Errorf("read csv header: %w", err)
	}

	normalizedHeader := normalizeRecord(header)
	if !slices.Equal(normalizedHeader, productionEntryHeader) {
		return nil, fmt.Errorf("expected header %v, got %v: %w", productionEntryHeader, normalizedHeader, ErrInvalidHeader)
	}
	reader.FieldsPerRecord = len(productionEntryHeader)

	return &ProductionEntryReader{
		reader:    reader,
		rowNumber: 1,
	}, nil
}

// Read returns the next raw production-entry row. It returns io.EOF when the stream is exhausted.
func (r *ProductionEntryReader) Read() (ProductionEntryRow, error) {
	record, err := r.reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return ProductionEntryRow{}, io.EOF
		}
		return ProductionEntryRow{}, fmt.Errorf("read csv record: %w", err)
	}

	r.rowNumber++
	if len(record) != len(productionEntryHeader) {
		return ProductionEntryRow{}, fmt.Errorf("row %d has %d fields, expected %d: %w", r.rowNumber, len(record), len(productionEntryHeader), ErrInvalidRecord)
	}

	return ProductionEntryRow{
		RowNumber:   r.rowNumber,
		EmployeeID:  strings.TrimSpace(record[0]),
		ProductSKU:  strings.TrimSpace(record[1]),
		Quantity:    strings.TrimSpace(record[2]),
		Workstation: strings.TrimSpace(record[3]),
		Timestamp:   strings.TrimSpace(record[4]),
		Comment:     strings.TrimSpace(record[5]),
	}, nil
}

func normalizeRecord(record []string) []string {
	normalized := make([]string, len(record))
	for i, field := range record {
		normalized[i] = strings.ToLower(strings.TrimSpace(field))
	}
	return normalized
}
