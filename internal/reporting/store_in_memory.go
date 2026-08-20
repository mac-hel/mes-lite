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
	mu                  sync.RWMutex
	dailyProductionRows []DailyProductionRow
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
