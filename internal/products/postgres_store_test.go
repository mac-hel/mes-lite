package products

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

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if _, err := pool.Exec(ctx, "DELETE FROM production_entries"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM products"); err != nil {
		t.Fatal(err)
	}
	return NewPostgresStore(pool)
}

func TestPostgresStore_SaveListUpdateAndConflict(t *testing.T) {
	store := testPostgresStore(t)
	prod, err := NewProduct("VX-100", "Ventilation Unit", "piece", CategoryVentilation)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), prod); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindBySKU(t.Context(), prod.SKU)
	if err != nil {
		t.Fatal(err)
	}
	if got != prod {
		t.Fatalf("expected %#v, got %#v", prod, got)
	}

	stale := prod
	if err := prod.UpdateDetails("Filter Set", "set", CategoryFilter); err != nil {
		t.Fatal(err)
	}
	updated, err := store.Update(t.Context(), prod)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != 2 {
		t.Fatalf("expected version 2, got %d", updated.Version)
	}
	if _, err := store.Update(t.Context(), stale); !errors.Is(err, ErrVersionConflict) {
		t.Fatalf("expected ErrVersionConflict, got %v", err)
	}
}
