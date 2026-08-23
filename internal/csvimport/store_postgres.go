package csvimport

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mac-hel/mes-lite/internal/postgres"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/production/productiondb"
)

// BatchError identifies which validated import row caused a batch persistence failure.
type BatchError struct {
	RowNumber int
	Err       error
}

// Error returns a human-readable batch persistence error.
func (e BatchError) Error() string {
	return fmt.Sprintf("row %d: %v", e.RowNumber, e.Err)
}

// Unwrap returns the underlying persistence error.
func (e BatchError) Unwrap() error {
	return e.Err
}

// NewPostgresStore creates a PostgreSQL-backed CSV import store.
func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// PostgresStore persists validated CSV import records to PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

// SaveBatch inserts all records in one transaction. Any failed row rolls back the whole batch.
func (s *PostgresStore) SaveBatch(ctx context.Context, records []ProductionEntryRecord) ([]production.Entry, error) {
	if len(records) == 0 {
		return []production.Entry{}, nil
	}

	entries := make([]production.Entry, 0, len(records))
	for _, record := range records {
		entry, err := entryFromRecord(record)
		if err != nil {
			return nil, BatchError{RowNumber: record.RowNumber, Err: err}
		}
		entries = append(entries, entry)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin csv import batch: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := productiondb.New(tx)
	for i, entry := range entries {
		if err := createEntry(ctx, queries, entry); err != nil {
			return nil, BatchError{RowNumber: records[i].RowNumber, Err: mapPostgresError(entry.ID, err)}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit csv import batch: %w", err)
	}

	return entries, nil
}

func entryFromRecord(record ProductionEntryRecord) (production.Entry, error) {
	id, err := production.NewEntryID()
	if err != nil {
		return production.Entry{}, err
	}

	return production.NewEntry(
		id,
		record.EmployeeID,
		record.ProductSKU,
		record.Quantity,
		record.Workstation,
		record.Timestamp,
		record.Comment,
	)
}

func createEntry(ctx context.Context, queries *productiondb.Queries, entry production.Entry) error {
	id, err := uuidFromString(entry.ID)
	if err != nil {
		return err
	}

	_, err = queries.CreateEntry(ctx, productiondb.CreateEntryParams{
		ID:          id,
		RequestID:   entry.RequestID,
		EmployeeID:  entry.EmployeeID,
		ProductSku:  entry.ProductSKU,
		Quantity:    int32(entry.Quantity),
		Workstation: entry.Workstation,
		OccurredAt:  pgtype.Timestamptz{Time: entry.Timestamp, Valid: true},
		Comment:     entry.Comment,
	})
	return err
}

func uuidFromString(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return pgtype.UUID{}, fmt.Errorf("id must be a UUID: %w", production.ErrInvalidEntry)
	}
	return uuid, nil
}

func mapPostgresError(id string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("production entry %q: %w", id, production.ErrNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			return fmt.Errorf("production entry %q: %w", id, production.ErrAlreadyExists)
		case postgres.CheckViolation, postgres.NotNullViolation, postgres.InvalidTextValue, postgres.ForeignKeyViolation:
			return fmt.Errorf("production entry %q: %w", id, production.ErrInvalidEntry)
		}
	}

	return err
}
