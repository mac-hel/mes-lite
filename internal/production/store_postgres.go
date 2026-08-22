package production

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/mac-hel/mes-lite/internal/postgres"
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

func mapPostgresError(id string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			return fmt.Errorf("production entry %q: %w", id, ErrAlreadyExists)
		case postgres.CheckViolation, postgres.NotNullViolation, postgres.InvalidTextValue, postgres.ForeignKeyViolation:
			return fmt.Errorf("production entry %q: %w", id, ErrInvalidEntry)
		}
	}

	return err
}

func entryFromDB(entry productiondb.ProductionEntry) (Entry, error) {
	return NewEntry(
		uuidString(entry.ID),
		entry.EmployeeID,
		entry.ProductSku,
		int(entry.Quantity),
		entry.Workstation,
		entry.OccurredAt.Time,
		entry.Comment,
	)
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
