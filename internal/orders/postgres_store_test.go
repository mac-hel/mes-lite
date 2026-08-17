package orders

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

	if _, err := pool.Exec(ctx, "DELETE FROM production_order_assignments"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM production_order_lines"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM production_orders"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM production_entries"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM employees"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM products"); err != nil {
		t.Fatal(err)
	}
	seedOrderReferences(t, ctx, pool)

	return NewPostgresStore(pool)
}

func TestPostgresStore_SaveAndFindByID(t *testing.T) {
	store := testPostgresStore(t)
	order := mustOrderWithLines(t,
		mustOrderLine(t, "shaft-1", 2),
		mustOrderLine(t, "filter-1", 4),
	)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}

	if err := store.Save(t.Context(), order); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindByID(t.Context(), order.ID())
	if err != nil {
		t.Fatal(err)
	}
	assertOrdersEqual(t, order, got)
}

func TestPostgresStore_SaveDuplicate(t *testing.T) {
	store := testPostgresStore(t)
	order := mustOrderWithLines(t, mustOrderLine(t, "shaft-1", 2))

	if err := store.Save(t.Context(), order); err != nil {
		t.Fatal(err)
	}

	err := store.Save(t.Context(), order)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("Save() error = %v, want ErrAlreadyExists", err)
	}
}

func TestPostgresStore_FindByID_NotFound(t *testing.T) {
	store := testPostgresStore(t)

	_, err := store.FindByID(t.Context(), "missing-order")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("FindByID() error = %v, want ErrNotFound", err)
	}
}

func TestPostgresStore_SaveMissingProductReferenceRollsBack(t *testing.T) {
	store := testPostgresStore(t)
	order := mustOrderWithLines(t, mustOrderLine(t, "missing-product", 2))

	err := store.Save(t.Context(), order)
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("Save() error = %v, want ErrInvalidOrder", err)
	}
	_, err = store.FindByID(t.Context(), order.ID())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected failed save to roll back order row, got %v", err)
	}
}

func TestPostgresStore_SaveMissingEmployeeReferenceRollsBack(t *testing.T) {
	store := testPostgresStore(t)
	order := mustOrderWithLines(t, mustOrderLine(t, "shaft-1", 2))
	if err := order.AssignEmployee("missing-employee", time.Now()); err != nil {
		t.Fatal(err)
	}

	err := store.Save(t.Context(), order)
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("Save() error = %v, want ErrInvalidOrder", err)
	}
	_, err = store.FindByID(t.Context(), order.ID())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected failed save to roll back order row, got %v", err)
	}
}

func seedOrderReferences(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO products (sku, name, category, unit, is_active, version)
		VALUES
			('shaft-1', 'Shaft', 0, 'piece', true, 1),
			('filter-1', 'Filter', 1, 'piece', true, 1)
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO employees (id, first_name, last_name, email, is_active, version)
		VALUES ('emp-1', 'Ana', 'Worker', 'ana@example.com', true, 1)
	`); err != nil {
		t.Fatal(err)
	}
}

func assertOrdersEqual(t *testing.T, want, got Order) {
	t.Helper()
	if got.ID() != want.ID() {
		t.Fatalf("ID = %q, want %q", got.ID(), want.ID())
	}
	if got.Status() != want.Status() {
		t.Fatalf("Status = %q, want %q", got.Status(), want.Status())
	}
	if !got.CreatedAt().Equal(want.CreatedAt()) {
		t.Fatalf("CreatedAt = %s, want %s", got.CreatedAt(), want.CreatedAt())
	}
	if !got.UpdatedAt().Equal(want.UpdatedAt()) {
		t.Fatalf("UpdatedAt = %s, want %s", got.UpdatedAt(), want.UpdatedAt())
	}
	assertOrderLinesEqual(t, want.Lines(), got.Lines())
	assertStringsEqual(t, want.AssignedEmployees(), got.AssignedEmployees())
}

func assertOrderLinesEqual(t *testing.T, want, got []OrderLine) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("line count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].ProductSKU() != want[i].ProductSKU() {
			t.Fatalf("line[%d].ProductSKU = %q, want %q", i, got[i].ProductSKU(), want[i].ProductSKU())
		}
		if got[i].PlannedQuantity() != want[i].PlannedQuantity() {
			t.Fatalf("line[%d].PlannedQuantity = %d, want %d", i, got[i].PlannedQuantity(), want[i].PlannedQuantity())
		}
	}
}

func assertStringsEqual(t *testing.T, want, got []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("string count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("string[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func mustOrderWithLines(t *testing.T, lines ...OrderLine) Order {
	t.Helper()
	order, err := NewOrder("order-1", lines, time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	return order
}
