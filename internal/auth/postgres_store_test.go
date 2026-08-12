package auth

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
	if _, err := pool.Exec(ctx, "DELETE FROM auth_users"); err != nil {
		t.Fatal(err)
	}
	return NewPostgresStore(pool)
}

func TestPostgresStoreSaveAndFindByEmail(t *testing.T) {
	store := testPostgresStore(t)
	user, err := NewUser("user-1", "ADMIN@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindByEmail(t.Context(), "admin@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != user.ID || got.Email != "admin@example.com" || got.Role != RoleAdmin || !got.IsActive {
		t.Fatalf("unexpected user: %+v", got)
	}
	if !got.VerifyPassword("secret") {
		t.Fatal("expected persisted password hash to verify")
	}
}

func TestPostgresStoreRejectsDuplicateEmail(t *testing.T) {
	store := testPostgresStore(t)
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}
	duplicate, err := NewUser("user-2", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(t.Context(), duplicate); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestPostgresStoreFindByEmailNotFound(t *testing.T) {
	store := testPostgresStore(t)

	_, err := store.FindByEmail(t.Context(), "missing@example.com")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresStoreEnsureBootstrapAdminIsIdempotent(t *testing.T) {
	store := testPostgresStore(t)
	user, err := NewUser("bootstrap-admin", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.EnsureBootstrapAdmin(t.Context(), user); err != nil {
		t.Fatal(err)
	}

	changed, err := NewUser("bootstrap-admin", "admin@example.com", "changed", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.EnsureBootstrapAdmin(t.Context(), changed); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindByEmail(t.Context(), "admin@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !got.VerifyPassword("secret") {
		t.Fatal("expected existing bootstrap password to remain unchanged")
	}
}
