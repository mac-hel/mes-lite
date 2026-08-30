package csvimport

import (
	"context"
	"errors"
	"io"

	"github.com/mac-hel/mes-lite/internal/production"
)

const (
	defaultImportBatchSize = 500
	maxReportedErrors      = 1000
)

// ImportSummary describes the outcome of a production-entry CSV import.
type ImportSummary struct {
	TotalRows         int           `json:"totalRows"`
	ValidRows         int           `json:"validRows"`
	InvalidRows       int           `json:"invalidRows"`
	ImportedRows      int           `json:"importedRows"`
	Errors            []ImportError `json:"errors"`
	ErrorLimitReached bool          `json:"errorLimitReached"`
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
	return s.importProductionEntries(ctx, input, defaultImportBatchSize)
}

func (s *Service) importProductionEntries(ctx context.Context, input io.Reader, batchSize int) (ImportSummary, error) {
	if batchSize <= 0 {
		batchSize = defaultImportBatchSize
	}

	reader, err := NewProductionEntryReader(input)
	if err != nil {
		return ImportSummary{}, err
	}

	summary := ImportSummary{Errors: []ImportError{}}
	batch := make([]ProductionEntryRecord, 0, batchSize)

	for {
		if err := ctx.Err(); err != nil {
			return ImportSummary{}, err
		}

		raw, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return s.saveImportBatch(ctx, summary, batch)
			}
			return ImportSummary{}, err
		}

		summary.TotalRows++
		record, rowErrors := validateProductionEntryRow(raw)
		if len(rowErrors) > 0 {
			summary.InvalidRows++
			summary.addErrors(importErrorsFromRowErrors(rowErrors))
			continue
		}

		summary.ValidRows++
		batch = append(batch, record)
		if len(batch) < batchSize {
			continue
		}

		var saveErr error
		summary, saveErr = s.saveImportBatch(ctx, summary, batch)
		if saveErr != nil {
			return ImportSummary{}, saveErr
		}
		batch = batch[:0]
	}
}

func (s *Service) saveImportBatch(ctx context.Context, summary ImportSummary, batch []ProductionEntryRecord) (ImportSummary, error) {
	if len(batch) == 0 {
		return summary, nil
	}

	entries, err := s.store.SaveBatch(ctx, batch)
	if err != nil {
		if isRowPersistenceError(err) {
			return s.saveImportRecordsIndividually(ctx, summary, batch)
		}
		return ImportSummary{}, err
	}

	summary.ImportedRows += len(entries)
	return summary, nil
}

func (s *Service) saveImportRecordsIndividually(ctx context.Context, summary ImportSummary, records []ProductionEntryRecord) (ImportSummary, error) {
	for _, record := range records {
		entries, err := s.store.SaveBatch(ctx, []ProductionEntryRecord{record})
		if err == nil {
			summary.ImportedRows += len(entries)
			continue
		}

		if !isRowPersistenceError(err) {
			return ImportSummary{}, err
		}

		summary.ValidRows--
		summary.InvalidRows++
		summary.addErrors([]ImportError{{
			RowNumber: record.RowNumber,
			Message:   rowPersistenceMessage(err),
		}})
	}

	return summary, nil
}

func isRowPersistenceError(err error) bool {
	var batchErr BatchError
	return errors.As(err, &batchErr) || errors.Is(err, production.ErrInvalidEntry) || errors.Is(err, production.ErrAlreadyExists)
}

func rowPersistenceMessage(err error) string {
	var batchErr BatchError
	if errors.As(err, &batchErr) {
		return batchErr.Err.Error()
	}
	return err.Error()
}

func (s *ImportSummary) addErrors(importErrors []ImportError) {
	for _, importError := range importErrors {
		if len(s.Errors) >= maxReportedErrors {
			s.ErrorLimitReached = true
			return
		}
		s.Errors = append(s.Errors, importError)
	}
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
