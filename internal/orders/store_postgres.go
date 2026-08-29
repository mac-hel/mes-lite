package orders

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mac-hel/mes-lite/internal/orders/ordersdb"
	"github.com/mac-hel/mes-lite/internal/platform/postgres"
)

// NewPostgresStore creates a PostgreSQL-backed [Store].
func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db, queries: ordersdb.New(db)}
}

// PostgresStore stores production orders in PostgreSQL through sqlc-generated queries.
type PostgresStore struct {
	db      *pgxpool.Pool
	queries *ordersdb.Queries
}

// Save stores a complete production order aggregate atomically.
func (s *PostgresStore) Save(ctx context.Context, order Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := s.queries.WithTx(tx)
	if _, err := queries.CreateOrder(ctx, ordersdb.CreateOrderParams{
		ID:        order.ID(),
		Status:    string(order.Status()),
		Version:   int32(order.Version()),
		CreatedAt: timestamptz(order.CreatedAt()),
		UpdatedAt: timestamptz(order.UpdatedAt()),
	}); err != nil {
		return mapPostgresError(order.ID(), err)
	}

	for _, line := range order.Lines().Values() {
		if line.PlannedQuantity() > math.MaxInt32 {
			return fmt.Errorf("production order %q: %w", order.ID(), ErrInvalidOrder)
		}
		if err := queries.CreateOrderLine(ctx, ordersdb.CreateOrderLineParams{
			OrderID:         order.ID(),
			ProductSku:      line.ProductSKU(),
			PlannedQuantity: int32(line.PlannedQuantity()),
		}); err != nil {
			return mapPostgresError(order.ID(), err)
		}
	}

	for _, employeeID := range order.AssignedEmployees() {
		if err := queries.CreateOrderAssignment(ctx, ordersdb.CreateOrderAssignmentParams{
			OrderID:    order.ID(),
			EmployeeID: employeeID,
		}); err != nil {
			return mapPostgresError(order.ID(), err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Update stores mutable production order state atomically and increments its version.
func (s *PostgresStore) Update(ctx context.Context, order Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := s.queries.WithTx(tx)
	if _, err := queries.UpdateOrder(ctx, ordersdb.UpdateOrderParams{
		ID:        order.ID(),
		Status:    string(order.Status()),
		UpdatedAt: timestamptz(order.UpdatedAt()),
		Version:   int32(order.Version()),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return s.updateNoRowsError(ctx, order.ID(), order.Version())
		}
		return mapPostgresError(order.ID(), err)
	}
	if err := queries.DeleteOrderAssignments(ctx, order.ID()); err != nil {
		return err
	}
	for _, employeeID := range order.AssignedEmployees() {
		if err := queries.CreateOrderAssignment(ctx, ordersdb.CreateOrderAssignmentParams{
			OrderID:    order.ID(),
			EmployeeID: employeeID,
		}); err != nil {
			return mapPostgresError(order.ID(), err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// FindByID looks up a production order by ID in PostgreSQL.
func (s *PostgresStore) FindByID(ctx context.Context, id string) (Order, error) {
	row, err := s.queries.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Order{}, fmt.Errorf("production order %q: %w", id, ErrNotFound)
		}
		return Order{}, err
	}

	lineRows, err := s.queries.ListOrderLines(ctx, id)
	if err != nil {
		return Order{}, err
	}
	lines := make([]OrderLine, 0, len(lineRows))
	for _, lineRow := range lineRows {
		line, err := NewOrderLine(lineRow.ProductSku, int(lineRow.PlannedQuantity))
		if err != nil {
			return Order{}, err
		}
		lines = append(lines, line)
	}

	assignmentRows, err := s.queries.ListOrderAssignments(ctx, id)
	if err != nil {
		return Order{}, err
	}
	assignedEmployees := make([]string, 0, len(assignmentRows))
	for _, assignmentRow := range assignmentRows {
		assignedEmployees = append(assignedEmployees, assignmentRow.EmployeeID)
	}

	orderLines, err := NewOrderLines(lines...)
	if err != nil {
		return Order{}, err
	}

	return RestoreOrder(row.ID, orderLines, Status(row.Status), assignedEmployees, int(row.Version), row.CreatedAt.Time, row.UpdatedAt.Time)
}

func (s *PostgresStore) updateNoRowsError(ctx context.Context, id string, version int) error {
	_, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return fmt.Errorf("production order %q version %d: %w", id, version, ErrVersionConflict)
}

func timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t.UTC(), Valid: true}
}

func mapPostgresError(id string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			return fmt.Errorf("production order %q: %w", id, ErrAlreadyExists)
		case postgres.CheckViolation, postgres.NotNullViolation, postgres.ForeignKeyViolation:
			return fmt.Errorf("production order %q: %w", id, ErrInvalidOrder)
		}
	}

	return err
}
