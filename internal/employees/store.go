package employees

import (
	"context"
	"fmt"
)

// Store defines the persistence operations for Employee data.
type Store interface {
	Save(ctx context.Context, emp Employee) error
	FindByID(ctx context.Context, id string) (Employee, error)
	List(ctx context.Context) ([]Employee, error)
	Update(ctx context.Context, emp Employee) error
}

// NewInMemoryStore creates a map-based in-memory [InMemoryStore].
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		employees: make(map[string]Employee),
	}
}

// InMemoryStore is a map-based in-memory implementation of [Store].
type InMemoryStore struct {
	employees map[string]Employee
}

// Save stores an employee keyed by ID.
func (s *InMemoryStore) Save(_ context.Context, emp Employee) error {
	s.employees[emp.ID] = emp
	return nil
}

// FindByID looks up an employee by ID. Returns [ErrNotFound] if not found.
func (s *InMemoryStore) FindByID(_ context.Context, id string) (Employee, error) {
	emp, ok := s.employees[id]
	if !ok {
		return Employee{}, fmt.Errorf("employee %q: %w", id, ErrNotFound)
	}
	return emp, nil
}

// List returns all employees.
func (s *InMemoryStore) List(_ context.Context) ([]Employee, error) {
	emps := make([]Employee, 0, len(s.employees))
	for _, emp := range s.employees {
		emps = append(emps, emp)
	}
	return emps, nil
}

// Update replaces the employee at the given ID. Returns [ErrNotFound] if not found.
func (s *InMemoryStore) Update(_ context.Context, emp Employee) error {
	if _, ok := s.employees[emp.ID]; !ok {
		return fmt.Errorf("employee %q: %w", emp.ID, ErrNotFound)
	}
	s.employees[emp.ID] = emp
	return nil
}
