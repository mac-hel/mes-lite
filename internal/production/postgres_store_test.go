package production

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
)

func testPostgresStore(t *testing.T) *PostgresStore {
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
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Logf("close migration db: %v", err)
		}
	})

	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("PostgreSQL is not available: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.ExecContext(ctx, "SELECT pg_advisory_lock(5001)"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := sqlDB.ExecContext(context.Background(), "SELECT pg_advisory_unlock(5001)"); err != nil {
			t.Logf("unlock migration lock: %v", err)
		}
	})
	if err := goose.Up(sqlDB, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if _, err := pool.Exec(ctx, "DELETE FROM production_entries"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM employees"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM products"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO employees (id, first_name, last_name, email, is_active, version)
		VALUES ('emp-1', 'Ana', 'Worker', 'ana@example.com', true, 1)
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO products (sku, name, category, unit, is_active, version)
		VALUES ('sku-1', 'Ventilation Unit', 0, 'piece', true, 1)
	`); err != nil {
		t.Fatal(err)
	}

	return NewPostgresStore(pool)
}

func TestPostgresStore_SaveAndFindByID(t *testing.T) {
	store := testPostgresStore(t)
	id, err := NewEntryID()
	if err != nil {
		t.Fatal(err)
	}
	entry, err := NewEntry(id, "emp-1", "sku-1", 12, "ws-1", time.Date(2026, 8, 8, 10, 30, 0, 0, time.UTC), "batch finished")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(t.Context(), entry); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindByID(t.Context(), entry.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got != entry {
		t.Errorf("expected %#v, got %#v", entry, got)
	}
}

func TestPostgresStore_SaveDuplicate(t *testing.T) {
	store := testPostgresStore(t)
	id, err := NewEntryID()
	if err != nil {
		t.Fatal(err)
	}
	entry, err := NewEntry(id, "emp-1", "sku-1", 12, "ws-1", time.Now(), "")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(t.Context(), entry); err != nil {
		t.Fatal(err)
	}

	err = store.Save(t.Context(), entry)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestPostgresStore_FindByID_NotFound(t *testing.T) {
	store := testPostgresStore(t)
	id, err := NewEntryID()
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.FindByID(t.Context(), id)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresStore_FindByID_InvalidID(t *testing.T) {
	store := testPostgresStore(t)

	_, err := store.FindByID(t.Context(), "not-a-uuid")
	if !errors.Is(err, ErrInvalidEntry) {
		t.Fatalf("expected ErrInvalidEntry, got %v", err)
	}
}

func TestPostgresStore_SaveMissingReference(t *testing.T) {
	store := testPostgresStore(t)
	id, err := NewEntryID()
	if err != nil {
		t.Fatal(err)
	}
	entry, err := NewEntry(id, "missing", "sku-1", 12, "ws-1", time.Now(), "")
	if err != nil {
		t.Fatal(err)
	}

	err = store.Save(t.Context(), entry)
	if !errors.Is(err, ErrInvalidEntry) {
		t.Fatalf("expected ErrInvalidEntry, got %v", err)
	}
	_, err = store.FindByID(t.Context(), entry.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected failed insert to leave no row, got %v", err)
	}
}

func TestPostgresStore_ListFiltersAndPaginates(t *testing.T) {
	store := testPostgresStore(t)
	entries := []Entry{
		mustProductionEntry(t, "00000000-0000-4000-8000-000000000011", "emp-1", "sku-1", 3, "assembly-1", "2026-08-08T10:30:00Z"),
		mustProductionEntry(t, "00000000-0000-4000-8000-000000000012", "emp-1", "sku-1", 5, "assembly-2", "2026-08-09T10:30:00Z"),
		mustProductionEntry(t, "00000000-0000-4000-8000-000000000013", "emp-1", "sku-1", 7, "paint", "2026-08-09T11:30:00Z"),
	}
	for _, entry := range entries {
		if err := store.Save(t.Context(), entry); err != nil {
			t.Fatal(err)
		}
	}

	got, err := store.List(t.Context(), ListOptions{
		ProductSKU:  "sku-1",
		Workstation: "assembly",
		From:        time.Date(2026, 8, 8, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC),
		Limit:       1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("expected one entry, got %d", len(got))
	}
	if got[0].ID != entries[1].ID {
		t.Fatalf("expected newest matching entry %q, got %q", entries[1].ID, got[0].ID)
	}
}
