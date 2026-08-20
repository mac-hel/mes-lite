package reporting

import (
	"context"
	"sync"
	"time"
)

// NewInMemoryStore creates an in-memory reporting store for fast HTTP tests.
func NewInMemoryStore(rows []DailyProductionRow) *InMemoryStore {
	copied := append([]DailyProductionRow(nil), rows...)
	return &InMemoryStore{dailyProductionRows: copied}
}

// InMemoryStore stores reporting rows in memory for tests.
type InMemoryStore struct {
	mu                       sync.RWMutex
	dailyProductionRows      []DailyProductionRow
	employeeProductivityRows []EmployeeProductivityRow
}

// NewInMemoryStoreWithReports creates an in-memory store with all supported report rows.
func NewInMemoryStoreWithReports(dailyRows []DailyProductionRow, employeeRows []EmployeeProductivityRow) *InMemoryStore {
	return &InMemoryStore{
		dailyProductionRows:      append([]DailyProductionRow(nil), dailyRows...),
		employeeProductivityRows: append([]EmployeeProductivityRow(nil), employeeRows...),
	}
}

// DailyProduction returns stored daily production rows whose day falls inside the requested range.
func (s *InMemoryStore) DailyProduction(_ context.Context, from, to time.Time) ([]DailyProductionRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	rows := make([]DailyProductionRow, 0, len(s.dailyProductionRows))
	from = from.UTC()
	to = to.UTC()
	for _, row := range s.dailyProductionRows {
		day := row.Day.UTC()
		if !day.Before(from) && day.Before(to) {
			row.Day = day
			rows = append(rows, row)
		}
	}

	return rows, nil
}

// EmployeeProductivity returns stored employee productivity rows for a valid range.
func (s *InMemoryStore) EmployeeProductivity(_ context.Context, from, to time.Time) ([]EmployeeProductivityRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]EmployeeProductivityRow(nil), s.employeeProductivityRows...), nil
}
