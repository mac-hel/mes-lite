package production

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/mac-hel/mes-lite/internal/config"
)

func testPostgresStore(t *testing.T) *PostgresStore {
	t.Helper()

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
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
	if err := goose.Up(sqlDB, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	if _, err := pool.Exec(ctx, "DELETE FROM production_entries"); err != nil {
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
	entry := NewEntry(id, "emp-1", "sku-1", 12, "ws-1", time.Date(2026, 8, 8, 10, 30, 0, 0, time.UTC), "batch finished")

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
	entry := NewEntry(id, "emp-1", "sku-1", 12, "ws-1", time.Now(), "")

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
