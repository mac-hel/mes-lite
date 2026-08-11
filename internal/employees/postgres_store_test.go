package employees

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
	if _, err := pool.Exec(ctx, "DELETE FROM production_entries"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM employees"); err != nil {
		t.Fatal(err)
	}
	return NewPostgresStore(pool)
}

func TestPostgresStore_SaveListUpdateAndConflict(t *testing.T) {
	store := testPostgresStore(t)
	emp, err := NewEmployee("emp-1", "John", "Worker", "john@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), emp); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindByID(t.Context(), emp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != emp {
		t.Fatalf("expected %#v, got %#v", emp, got)
	}

	stale := emp
	if err := emp.UpdateDetails("Jane", "Worker", "jane@example.com"); err != nil {
		t.Fatal(err)
	}
	updated, err := store.Update(t.Context(), emp)
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

func TestPostgresStore_UpdateConcurrentVersionConflict(t *testing.T) {
	store := testPostgresStore(t)
	emp, err := NewEmployee("emp-1", "John", "Worker", "john@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), emp); err != nil {
		t.Fatal(err)
	}

	first := emp
	_ = first.UpdateDetails("Jane", "Worker", "jane@example.com")
	second := emp
	_ = second.UpdateDetails("Ana", "Worker", "ana@example.com")

	errs := make(chan error, 2)
	go func() { _, err := store.Update(t.Context(), first); errs <- err }()
	go func() { _, err := store.Update(t.Context(), second); errs <- err }()

	successes, conflicts := 0, 0
	for range 2 {
		err := <-errs
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrVersionConflict):
			conflicts++
		default:
			t.Fatalf("unexpected update error: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("expected 1 success and 1 conflict, got %d successes and %d conflicts", successes, conflicts)
	}
}
