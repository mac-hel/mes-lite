package employees

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// NewInMemoryStore creates a map-based in-memory [InMemoryStore].
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{employees: make(map[string]Employee)}
}

// InMemoryStore is a map-based in-memory implementation of [Store].
type InMemoryStore struct {
	employees map[string]Employee
}

// Save stores an employee keyed by ID.
func (s *InMemoryStore) Save(_ context.Context, emp Employee) error {
	if err := emp.Validate(); err != nil {
		return err
	}
	if _, ok := s.employees[emp.ID]; ok {
		return fmt.Errorf("employee %q: %w", emp.ID, ErrAlreadyExists)
	}
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

// List returns employees matching the given options.
func (s *InMemoryStore) List(_ context.Context, opts ListOptions) ([]Employee, error) {
	opts, err := opts.normalize()
	if err != nil {
		return nil, err
	}
	emps := make([]Employee, 0, len(s.employees))
	for _, emp := range s.employees {
		if !matchesListOptions(emp, opts) {
			continue
		}
		emps = append(emps, emp)
	}
	sortEmployees(emps, opts.Sort)
	return paginate(emps, opts), nil
}

// Update replaces the employee at the given ID and increments its version.
func (s *InMemoryStore) Update(_ context.Context, emp Employee) (Employee, error) {
	if err := emp.Validate(); err != nil {
		return Employee{}, err
	}
	current, ok := s.employees[emp.ID]
	if !ok {
		return Employee{}, fmt.Errorf("employee %q: %w", emp.ID, ErrNotFound)
	}
	if current.Version != emp.Version {
		return Employee{}, fmt.Errorf("employee %q version %d: %w", emp.ID, emp.Version, ErrVersionConflict)
	}
	emp.Version++
	s.employees[emp.ID] = emp
	return emp, nil
}

func matchesListOptions(emp Employee, opts ListOptions) bool {
	if opts.Active != nil && emp.IsActive != *opts.Active {
		return false
	}
	if opts.Query == "" {
		return true
	}
	query := strings.ToLower(opts.Query)
	return strings.Contains(strings.ToLower(emp.ID), query) ||
		strings.Contains(strings.ToLower(emp.FirstName), query) ||
		strings.Contains(strings.ToLower(emp.LastName), query) ||
		strings.Contains(strings.ToLower(emp.Email), query)
}

func sortEmployees(emps []Employee, sortKey string) {
	sort.Slice(emps, func(i, j int) bool {
		a, b := emps[i], emps[j]
		switch sortKey {
		case "-id":
			return a.ID > b.ID
		case "name":
			if a.LastName == b.LastName {
				return a.FirstName < b.FirstName
			}
			return a.LastName < b.LastName
		case "-name":
			if a.LastName == b.LastName {
				return a.FirstName > b.FirstName
			}
			return a.LastName > b.LastName
		case "email":
			return a.Email < b.Email
		case "-email":
			return a.Email > b.Email
		default:
			return a.ID < b.ID
		}
	})
}

func paginate(emps []Employee, opts ListOptions) []Employee {
	if opts.Offset >= len(emps) {
		return []Employee{}
	}
	end := opts.Offset + opts.Limit
	if end > len(emps) {
		end = len(emps)
	}
	return emps[opts.Offset:end]
}
