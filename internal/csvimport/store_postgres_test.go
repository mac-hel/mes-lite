package csvimport

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/mac-hel/mes-lite/internal/production"
)

func TestPostgresStore_SaveBatchEmpty(t *testing.T) {
	store := &PostgresStore{}

	entries, err := store.SaveBatch(t.Context(), nil)
	if err != nil {
		t.Fatalf("SaveBatch() error = %v", err)
	}
	if entries == nil {
		t.Fatal("SaveBatch() entries = nil, want empty non-nil slice")
	}
	if len(entries) != 0 {
		t.Fatalf("SaveBatch() entries length = %d, want 0", len(entries))
	}
}

func TestPostgresStore_SaveBatchPersistsEntries(t *testing.T) {
	store, pool := testCSVImportPostgresStore(t)
	records := []ProductionEntryRecord{
		{
			RowNumber:   2,
			EmployeeID:  "emp-1",
			ProductSKU:  "sku-1",
			Quantity:    12,
			Workstation: "ws-1",
			Timestamp:   time.Date(2026, 8, 20, 10, 30, 0, 0, time.UTC),
			Comment:     "first",
		},
		{
			RowNumber:   3,
			EmployeeID:  "emp-2",
			ProductSKU:  "sku-2",
			Quantity:    7,
			Workstation: "ws-2",
			Timestamp:   time.Date(2026, 8, 20, 11, 30, 0, 0, time.UTC),
			Comment:     "second",
		},
	}

	entries, err := store.SaveBatch(t.Context(), records)
	if err != nil {
		t.Fatalf("SaveBatch() error = %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("SaveBatch() entries length = %d, want 2", len(entries))
	}
	for i, entry := range entries {
		if entry.ID == "" {
			t.Fatalf("entries[%d].ID is empty", i)
		}
		if entry.EmployeeID != records[i].EmployeeID || entry.ProductSKU != records[i].ProductSKU || entry.Quantity != records[i].Quantity || entry.Workstation != records[i].Workstation || !entry.Timestamp.Equal(records[i].Timestamp) || entry.Comment != records[i].Comment {
			t.Fatalf("entries[%d] = %+v, record = %+v", i, entry, records[i])
		}
	}

	var count int
	if err := pool.QueryRow(t.Context(), "SELECT count(*) FROM production_entries").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("persisted count = %d, want 2", count)
	}

	productionStore := production.NewPostgresStore(pool)
	for _, entry := range entries {
		got, err := productionStore.FindByID(t.Context(), entry.ID)
		if err != nil {
			t.Fatalf("FindByID(%q) error = %v", entry.ID, err)
		}
		if got != entry {
			t.Fatalf("FindByID(%q) = %+v, want %+v", entry.ID, got, entry)
		}
	}
}

func TestPostgresStore_SaveBatchRollsBackOnFailedRow(t *testing.T) {
	store, pool := testCSVImportPostgresStore(t)
	records := []ProductionEntryRecord{
		{
			RowNumber:   2,
			EmployeeID:  "emp-1",
			ProductSKU:  "sku-1",
			Quantity:    12,
			Workstation: "ws-1",
			Timestamp:   time.Date(2026, 8, 20, 10, 30, 0, 0, time.UTC),
		},
		{
			RowNumber:   3,
			EmployeeID:  "missing",
			ProductSKU:  "sku-1",
			Quantity:    7,
			Workstation: "ws-2",
			Timestamp:   time.Date(2026, 8, 20, 11, 30, 0, 0, time.UTC),
		},
	}

	_, err := store.SaveBatch(t.Context(), records)
	if !errors.Is(err, production.ErrInvalidEntry) {
		t.Fatalf("SaveBatch() error = %v, want production.ErrInvalidEntry", err)
	}
	var batchErr BatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("SaveBatch() error = %T, want BatchError", err)
	}
	if batchErr.RowNumber != 3 {
		t.Fatalf("BatchError.RowNumber = %d, want 3", batchErr.RowNumber)
	}

	var count int
	if err := pool.QueryRow(t.Context(), "SELECT count(*) FROM production_entries").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("persisted count after failed batch = %d, want 0", count)
	}
}

func testCSVImportPostgresStore(t *testing.T) (*PostgresStore, *pgxpool.Pool) {
	t.Helper()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("PostgreSQL is not available: %v", err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.ExecContext(ctx, "SELECT pg_advisory_lock(5001)"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = sqlDB.ExecContext(context.Background(), "SELECT pg_advisory_unlock(5001)") })
	if err := goose.Up(sqlDB, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	cleanCSVImportTables(t, ctx, pool)
	seedCSVImportReferences(t, ctx, pool)

	return NewPostgresStore(pool), pool
}

func cleanCSVImportTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	for _, query := range []string{
		"DELETE FROM production_order_assignments",
		"DELETE FROM production_order_lines",
		"DELETE FROM production_orders",
		"DELETE FROM production_entries",
		"DELETE FROM employees",
		"DELETE FROM products",
	} {
		if _, err := pool.Exec(ctx, query); err != nil {
			t.Fatal(err)
		}
	}
}

func seedCSVImportReferences(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO employees (id, first_name, last_name, email, is_active, version)
		VALUES
			('emp-1', 'Ana', 'Worker', 'ana@example.com', true, 1),
			('emp-2', 'Ben', 'Worker', 'ben@example.com', true, 1)
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO products (sku, name, category, unit, is_active, version)
		VALUES
			('sku-1', 'Ventilation Unit', 0, 'piece', true, 1),
			('sku-2', 'Filter', 1, 'piece', true, 1)
	`); err != nil {
		t.Fatal(err)
	}
}
