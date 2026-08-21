package csvimport

import (
	"context"
	"errors"
	"io"
)

// ImportSummary describes the outcome of a production-entry CSV import.
type ImportSummary struct {
	TotalRows    int           `json:"totalRows"`
	ValidRows    int           `json:"validRows"`
	InvalidRows  int           `json:"invalidRows"`
	ImportedRows int           `json:"importedRows"`
	Errors       []ImportError `json:"errors"`
}

// ImportError describes one row-level import problem returned to API clients.
type ImportError struct {
	RowNumber int    `json:"rowNumber"`
	Field     string `json:"field,omitempty"`
	Message   string `json:"message"`
}

// Service coordinates CSV import reading, validation and persistence.
type Service struct {
	store Store
}

// NewService creates a CSV import service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// ImportProductionEntries imports production entries from a CSV stream.
func (s *Service) ImportProductionEntries(ctx context.Context, input io.Reader) (ImportSummary, error) {
	reader, err := NewProductionEntryReader(input)
	if err != nil {
		return ImportSummary{}, err
	}

	validation, err := ValidateProductionEntries(reader)
	if err != nil {
		return ImportSummary{}, err
	}

	summary := ImportSummary{
		TotalRows:   len(validation.Records) + countInvalidRows(validation.Errors),
		ValidRows:   len(validation.Records),
		InvalidRows: countInvalidRows(validation.Errors),
		Errors:      importErrorsFromRowErrors(validation.Errors),
	}
	if len(validation.Records) == 0 {
		return summary, nil
	}

	entries, err := s.store.SaveBatch(ctx, validation.Records)
	if err != nil {
		var batchErr BatchError
		if errors.As(err, &batchErr) {
			summary.ValidRows = len(validation.Records) - 1
			summary.InvalidRows++
			summary.Errors = append(summary.Errors, ImportError{
				RowNumber: batchErr.RowNumber,
				Message:   batchErr.Err.Error(),
			})
			return summary, nil
		}
		return ImportSummary{}, err
	}

	summary.ImportedRows = len(entries)
	return summary, nil
}

func countInvalidRows(rowErrors []RowError) int {
	seen := make(map[int]struct{})
	for _, rowError := range rowErrors {
		seen[rowError.RowNumber] = struct{}{}
	}
	return len(seen)
}

func importErrorsFromRowErrors(rowErrors []RowError) []ImportError {
	if len(rowErrors) == 0 {
		return []ImportError{}
	}

	errors := make([]ImportError, 0, len(rowErrors))
	for _, rowError := range rowErrors {
		errors = append(errors, ImportError(rowError))
	}
	return errors
}
