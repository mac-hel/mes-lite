package production

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/mac-hel/mes-lite/internal/platform/postgres"
	"github.com/mac-hel/mes-lite/internal/production/productiondb"
)

// NewPostgresStore creates a PostgreSQL-backed [Store].
func NewPostgresStore(db productiondb.DBTX) *PostgresStore {
	return &PostgresStore{queries: productiondb.New(db)}
}

// PostgresStore stores production entries in PostgreSQL through sqlc-generated queries.
type PostgresStore struct {
	queries *productiondb.Queries
}

// Save stores a production entry in PostgreSQL.
func (s *PostgresStore) Save(ctx context.Context, entry Entry) error {
	if err := entry.Validate(); err != nil {
		return err
	}

	id, err := parseUUID(entry.ID)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateEntry(ctx, productiondb.CreateEntryParams{
		ID:          id,
		RequestID:   entry.RequestID,
		EmployeeID:  entry.EmployeeID,
		ProductSku:  entry.ProductSKU,
		Quantity:    int32(entry.Quantity),
		Workstation: entry.Workstation,
		OccurredAt:  pgtype.Timestamptz{Time: entry.Timestamp, Valid: true},
		Comment:     entry.Comment,
	})
	if err != nil {
		return mapPostgresError(entry.ID, err)
	}

	return nil
}

// FindByID looks up a production entry by ID in PostgreSQL.
func (s *PostgresStore) FindByID(ctx context.Context, id string) (Entry, error) {
	uuid, err := parseUUID(id)
	if err != nil {
		return Entry{}, err
	}

	entry, err := s.queries.GetEntry(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Entry{}, fmt.Errorf("production entry %q: %w", id, ErrNotFound)
		}
		return Entry{}, err
	}

	return entryFromDB(entry)
}

// FindByRequestID looks up a production entry by idempotency request ID in PostgreSQL.
func (s *PostgresStore) FindByRequestID(ctx context.Context, requestID string) (Entry, error) {
	entry, err := s.queries.GetEntryByRequestID(ctx, strings.TrimSpace(requestID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Entry{}, fmt.Errorf("production request id %q: %w", requestID, ErrNotFound)
		}
		return Entry{}, err
	}

	return entryFromDB(entry)
}

// List returns production entries matching review filters, newest first.
func (s *PostgresStore) List(ctx context.Context, opts ListOptions) ([]Entry, error) {
	opts, err := opts.normalize()
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.ListEntries(ctx, productiondb.ListEntriesParams{
		EmployeeID:  opts.EmployeeID,
		ProductSku:  opts.ProductSKU,
		Workstation: opts.Workstation,
		FromTime:    !opts.From.IsZero(),
		FromValue:   pgtype.Timestamptz{Time: opts.From, Valid: !opts.From.IsZero()},
		ToTime:      !opts.To.IsZero(),
		ToValue:     pgtype.Timestamptz{Time: opts.To, Valid: !opts.To.IsZero()},
		LimitValue:  int32(opts.Limit),
		OffsetValue: int32(opts.Offset),
	})
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(rows))
	for _, row := range rows {
		entry, err := entryFromDB(row)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// SaveCorrection stores an append-only production-entry correction in PostgreSQL.
func (s *PostgresStore) SaveCorrection(ctx context.Context, correction Correction) error {
	if err := correction.Validate(); err != nil {
		return err
	}
	id, err := parseUUID(correction.ID)
	if err != nil {
		return err
	}
	entryID, err := parseUUID(correction.EntryID)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateCorrection(ctx, productiondb.CreateCorrectionParams{
		ID:          id,
		EntryID:     entryID,
		ActorUserID: correction.ActorUserID,
		Reason:      correction.Reason,
		EmployeeID:  correction.EmployeeID,
		ProductSku:  correction.ProductSKU,
		Quantity:    int32(correction.Quantity),
		Workstation: correction.Workstation,
		OccurredAt:  pgtype.Timestamptz{Time: correction.Timestamp, Valid: true},
		Comment:     correction.Comment,
		CreatedAt:   pgtype.Timestamptz{Time: correction.CreatedAt, Valid: true},
	})
	if err != nil {
		return mapCorrectionPostgresError(correction.ID, err)
	}

	return nil
}

// ListCorrections returns correction history for a production entry from PostgreSQL.
func (s *PostgresStore) ListCorrections(ctx context.Context, entryID string) ([]Correction, error) {
	uuid, err := parseUUID(entryID)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.ListCorrections(ctx, uuid)
	if err != nil {
		return nil, err
	}
	corrections := make([]Correction, 0, len(rows))
	for _, row := range rows {
		correction, err := correctionFromDB(row)
		if err != nil {
			return nil, err
		}
		corrections = append(corrections, correction)
	}

	return corrections, nil
}

func mapPostgresError(id string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			if pgErr.ConstraintName == "production_entries_request_id_key" {
				return fmt.Errorf("production entry %q: %w", id, ErrRequestConflict)
			}
			return fmt.Errorf("production entry %q: %w", id, ErrAlreadyExists)
		case postgres.CheckViolation, postgres.NotNullViolation, postgres.InvalidTextValue, postgres.ForeignKeyViolation:
			return fmt.Errorf("production entry %q: %w", id, ErrInvalidEntry)
		}
	}

	return err
}

func mapCorrectionPostgresError(id string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			return fmt.Errorf("production correction %q: %w", id, ErrAlreadyExists)
		case postgres.ForeignKeyViolation:
			return fmt.Errorf("production correction %q: %w", id, ErrNotFound)
		case postgres.CheckViolation, postgres.NotNullViolation, postgres.InvalidTextValue:
			return fmt.Errorf("production correction %q: %w", id, ErrInvalidCorrection)
		}
	}

	return err
}

func entryFromDB(entry productiondb.ProductionEntry) (Entry, error) {
	return NewEntryWithRequestID(
		uuidString(entry.ID),
		entry.RequestID,
		entry.EmployeeID,
		entry.ProductSku,
		int(entry.Quantity),
		entry.Workstation,
		entry.OccurredAt.Time,
		entry.Comment,
	)
}

func correctionFromDB(correction productiondb.ProductionEntryCorrection) (Correction, error) {
	createdAt := timeFromPg(correction.CreatedAt)
	result, err := NewCorrection(
		uuidString(correction.ID),
		uuidString(correction.EntryID),
		correction.ActorUserID,
		correction.Reason,
		correction.EmployeeID,
		correction.ProductSku,
		int(correction.Quantity),
		correction.Workstation,
		correction.OccurredAt.Time,
		correction.Comment,
	)
	if err != nil {
		return Correction{}, err
	}
	result.CreatedAt = createdAt
	return result, nil
}

func timeFromPg(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time.UTC()
}

func parseUUID(id string) (pgtype.UUID, error) {
	compact := strings.ReplaceAll(strings.TrimSpace(id), "-", "")
	if len(compact) != 32 {
		return pgtype.UUID{}, fmt.Errorf("id must be a UUID: %w", ErrInvalidEntry)
	}

	var uuid pgtype.UUID
	if _, err := hex.Decode(uuid.Bytes[:], []byte(compact)); err != nil {
		return pgtype.UUID{}, fmt.Errorf("id must be a UUID: %w", ErrInvalidEntry)
	}
	uuid.Valid = true

	return uuid, nil
}

func uuidString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}

	encoded := hex.EncodeToString(uuid.Bytes[:])
	return encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:]
}
