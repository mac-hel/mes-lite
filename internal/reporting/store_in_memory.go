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
	productStatisticsRows    []ProductStatisticsRow
	dailyEmployeeRows        []DailyEmployeeProductionRow
	employeeProductRows      []EmployeeProductivityProductRow
}

// NewInMemoryStoreWithReports creates an in-memory store with all supported report rows.
func NewInMemoryStoreWithReports(dailyRows []DailyProductionRow, employeeRows []EmployeeProductivityRow, productRows []ProductStatisticsRow, dailyEmployeeRows []DailyEmployeeProductionRow, employeeProductRows []EmployeeProductivityProductRow) *InMemoryStore {
	return &InMemoryStore{
		dailyProductionRows:      append([]DailyProductionRow(nil), dailyRows...),
		employeeProductivityRows: append([]EmployeeProductivityRow(nil), employeeRows...),
		productStatisticsRows:    append([]ProductStatisticsRow(nil), productRows...),
		dailyEmployeeRows:        append([]DailyEmployeeProductionRow(nil), dailyEmployeeRows...),
		employeeProductRows:      append([]EmployeeProductivityProductRow(nil), employeeProductRows...),
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

// ProductStatistics returns stored product statistics rows for a valid range.
func (s *InMemoryStore) ProductStatistics(_ context.Context, from, to time.Time) ([]ProductStatisticsRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]ProductStatisticsRow(nil), s.productStatisticsRows...), nil
}

// DailyEmployeeProduction returns stored daily employee production rows for a valid range.
func (s *InMemoryStore) DailyEmployeeProduction(_ context.Context, from, to time.Time) ([]DailyEmployeeProductionRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]DailyEmployeeProductionRow(nil), s.dailyEmployeeRows...), nil
}

// EmployeeProductivityProducts returns stored employee product productivity rows for a valid range.
func (s *InMemoryStore) EmployeeProductivityProducts(_ context.Context, from, to time.Time) ([]EmployeeProductivityProductRow, error) {
	if err := validateRange(from, to); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]EmployeeProductivityProductRow(nil), s.employeeProductRows...), nil
}
