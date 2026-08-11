package employees

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mac-hel/mes-lite/internal/employees/employeesdb"
	"github.com/mac-hel/mes-lite/internal/postgres"
)

// NewPostgresStore creates a PostgreSQL-backed [Store].
func NewPostgresStore(db employeesdb.DBTX) *PostgresStore {
	return &PostgresStore{queries: employeesdb.New(db)}
}

// PostgresStore stores employees in PostgreSQL through sqlc-generated queries.
type PostgresStore struct {
	queries *employeesdb.Queries
}

// Save stores an employee keyed by ID.
func (s *PostgresStore) Save(ctx context.Context, emp Employee) error {
	if err := emp.Validate(); err != nil {
		return err
	}
	_, err := s.queries.CreateEmployee(ctx, employeesdb.CreateEmployeeParams{
		ID:        emp.ID,
		FirstName: emp.FirstName,
		LastName:  emp.LastName,
		Email:     emp.Email,
		IsActive:  emp.IsActive,
		Version:   int32(emp.Version),
	})
	if err != nil {
		return mapPostgresError(emp.ID, err)
	}
	return nil
}

// FindByID looks up an employee by ID. Returns [ErrNotFound] if not found.
func (s *PostgresStore) FindByID(ctx context.Context, id string) (Employee, error) {
	emp, err := s.queries.GetEmployee(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Employee{}, fmt.Errorf("employee %q: %w", id, ErrNotFound)
		}
		return Employee{}, err
	}
	return employeeFromGetDB(emp), nil
}

// List returns employees matching the given options.
func (s *PostgresStore) List(ctx context.Context, opts ListOptions) ([]Employee, error) {
	opts, err := opts.normalize()
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.ListEmployees(ctx, employeesdb.ListEmployeesParams{
		Query:       opts.Query,
		Active:      activeFilter(opts.Active),
		Sort:        opts.Sort,
		OffsetValue: int32(opts.Offset),
		LimitValue:  int32(opts.Limit),
	})
	if err != nil {
		return nil, err
	}
	emps := make([]Employee, 0, len(rows))
	for _, row := range rows {
		emps = append(emps, employeeFromListDB(row))
	}
	return emps, nil
}

// Update replaces the employee at the given ID and increments its version.
func (s *PostgresStore) Update(ctx context.Context, emp Employee) (Employee, error) {
	if err := emp.Validate(); err != nil {
		return Employee{}, err
	}
	row, err := s.queries.UpdateEmployee(ctx, employeesdb.UpdateEmployeeParams{
		ID:        emp.ID,
		FirstName: emp.FirstName,
		LastName:  emp.LastName,
		Email:     emp.Email,
		IsActive:  emp.IsActive,
		Version:   int32(emp.Version),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return s.updateNoRowsError(ctx, emp.ID, emp.Version)
		}
		return Employee{}, mapPostgresError(emp.ID, err)
	}
	return employeeFromUpdateDB(row), nil
}

func employeeFromGetDB(emp employeesdb.GetEmployeeRow) Employee {
	return Employee{ID: emp.ID, FirstName: emp.FirstName, LastName: emp.LastName, Email: emp.Email, IsActive: emp.IsActive, Version: int(emp.Version)}
}

func employeeFromListDB(emp employeesdb.ListEmployeesRow) Employee {
	return Employee{ID: emp.ID, FirstName: emp.FirstName, LastName: emp.LastName, Email: emp.Email, IsActive: emp.IsActive, Version: int(emp.Version)}
}

func employeeFromUpdateDB(emp employeesdb.UpdateEmployeeRow) Employee {
	return Employee{ID: emp.ID, FirstName: emp.FirstName, LastName: emp.LastName, Email: emp.Email, IsActive: emp.IsActive, Version: int(emp.Version)}
}

func (s *PostgresStore) updateNoRowsError(ctx context.Context, id string, version int) (Employee, error) {
	_, err := s.FindByID(ctx, id)
	if err != nil {
		return Employee{}, err
	}
	return Employee{}, fmt.Errorf("employee %q version %d: %w", id, version, ErrVersionConflict)
}

func mapPostgresError(id string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			return fmt.Errorf("employee %q: %w", id, ErrAlreadyExists)
		case postgres.CheckViolation, postgres.NotNullViolation:
			return fmt.Errorf("employee %q: %w", id, ErrInvalidEmployee)
		}
	}
	return err
}
