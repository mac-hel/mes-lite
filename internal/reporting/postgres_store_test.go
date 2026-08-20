package reporting

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

func testPostgresStore(t *testing.T) (*PostgresStore, *pgxpool.Pool) {
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

	cleanReportingTables(t, ctx, pool)
	seedReportingReferences(t, ctx, pool)

	return NewPostgresStore(pool), pool
}

func TestPostgresStore_DailyProduction(t *testing.T) {
	store, pool := testPostgresStore(t)
	ctx := t.Context()
	day1 := time.Date(2026, 8, 18, 10, 30, 0, 0, time.UTC)
	day2 := time.Date(2026, 8, 19, 7, 15, 0, 0, time.UTC)
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000001", "emp-1", "shaft-1", 3, day1)
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000002", "emp-2", "shaft-1", 4, day1.Add(2*time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000003", "emp-1", "filter-1", 5, day2)
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000004", "emp-1", "shaft-1", 9, day2.Add(48*time.Hour))

	got, err := store.DailyProduction(ctx, day1.Add(-time.Hour), day2.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	want := []DailyProductionRow{
		{Day: time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC), ProductSKU: "shaft-1", TotalQuantity: 7, EntryCount: 2},
		{Day: time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC), ProductSKU: "filter-1", TotalQuantity: 5, EntryCount: 1},
	}
	if len(got) != len(want) {
		t.Fatalf("len(DailyProduction()) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("DailyProduction()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestPostgresStore_DailyProductionInvalidRange(t *testing.T) {
	store, _ := testPostgresStore(t)
	now := time.Now()
	_, err := store.DailyProduction(t.Context(), now, now)
	if !errors.Is(err, ErrInvalidRange) {
		t.Fatalf("DailyProduction() error = %v, want ErrInvalidRange", err)
	}
}

func TestPostgresStore_EmployeeProductivity(t *testing.T) {
	store, pool := testPostgresStore(t)
	ctx := t.Context()
	from := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000011", "emp-1", "shaft-1", 3, from.Add(time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000012", "emp-1", "filter-1", 4, from.Add(2*time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000013", "emp-2", "shaft-1", 5, from.Add(3*time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000014", "emp-2", "shaft-1", 9, to.Add(time.Hour))

	got, err := store.EmployeeProductivity(ctx, from, to)
	if err != nil {
		t.Fatal(err)
	}

	want := []EmployeeProductivityRow{
		{EmployeeID: "emp-1", FirstName: "Ana", LastName: "Worker", TotalQuantity: 7, EntryCount: 2},
		{EmployeeID: "emp-2", FirstName: "Ben", LastName: "Worker", TotalQuantity: 5, EntryCount: 1},
	}
	if len(got) != len(want) {
		t.Fatalf("len(EmployeeProductivity()) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("EmployeeProductivity()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestPostgresStore_EmployeeProductivityInvalidRange(t *testing.T) {
	store, _ := testPostgresStore(t)
	now := time.Now()
	_, err := store.EmployeeProductivity(t.Context(), now, now)
	if !errors.Is(err, ErrInvalidRange) {
		t.Fatalf("EmployeeProductivity() error = %v, want ErrInvalidRange", err)
	}
}

func TestPostgresStore_ProductStatistics(t *testing.T) {
	store, pool := testPostgresStore(t)
	ctx := t.Context()
	from := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000021", "emp-1", "shaft-1", 3, from.Add(time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000022", "emp-2", "shaft-1", 4, from.Add(2*time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000023", "emp-1", "filter-1", 5, from.Add(3*time.Hour))
	insertReportingEntry(t, ctx, pool, "00000000-0000-4000-8000-000000000024", "emp-2", "filter-1", 9, to.Add(time.Hour))

	got, err := store.ProductStatistics(ctx, from, to)
	if err != nil {
		t.Fatal(err)
	}

	want := []ProductStatisticsRow{
		{ProductSKU: "shaft-1", ProductName: "Shaft", TotalQuantity: 7, EntryCount: 2, EmployeeCount: 2},
		{ProductSKU: "filter-1", ProductName: "Filter", TotalQuantity: 5, EntryCount: 1, EmployeeCount: 1},
	}
	if len(got) != len(want) {
		t.Fatalf("len(ProductStatistics()) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ProductStatistics()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestPostgresStore_ProductStatisticsInvalidRange(t *testing.T) {
	store, _ := testPostgresStore(t)
	now := time.Now()
	_, err := store.ProductStatistics(t.Context(), now, now)
	if !errors.Is(err, ErrInvalidRange) {
		t.Fatalf("ProductStatistics() error = %v, want ErrInvalidRange", err)
	}
}

func TestReportingIndexExists(t *testing.T) {
	_, pool := testPostgresStore(t)

	var exists bool
	if err := pool.QueryRow(t.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = 'public'
			  AND tablename = 'production_entries'
			  AND indexname = 'production_entries_reporting_idx'
		)
	`).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("production_entries_reporting_idx does not exist")
	}
}

func cleanReportingTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
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

func seedReportingReferences(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
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
		VALUES
			('emp-1', 'Ana', 'Worker', 'ana@example.com', true, 1),
			('emp-2', 'Ben', 'Worker', 'ben@example.com', true, 1)
	`); err != nil {
		t.Fatal(err)
	}
}

func insertReportingEntry(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id, employeeID, productSKU string, quantity int, occurredAt time.Time) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO production_entries (id, employee_id, product_sku, quantity, workstation, occurred_at, comment)
		VALUES ($1, $2, $3, $4, 'assembly', $5, '')
	`, id, employeeID, productSKU, quantity, occurredAt)
	if err != nil {
		t.Fatal(err)
	}
}
