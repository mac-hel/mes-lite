package csvimport

import (
	"context"
	"sync"

	"github.com/mac-hel/mes-lite/internal/production"
)

// InMemoryStore stores imported records in memory for fast handler and server tests.
type InMemoryStore struct {
	mu      sync.Mutex
	entries []production.Entry
	err     error
}

// NewInMemoryStore creates an empty in-memory CSV import store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

// NewFailingInMemoryStore creates an in-memory CSV import store that returns err from SaveBatch.
func NewFailingInMemoryStore(err error) *InMemoryStore {
	return &InMemoryStore{err: err}
}

// SaveBatch converts records to production entries and stores them in memory.
func (s *InMemoryStore) SaveBatch(_ context.Context, records []ProductionEntryRecord) ([]production.Entry, error) {
	if s.err != nil {
		return nil, s.err
	}
	if len(records) == 0 {
		return []production.Entry{}, nil
	}

	entries := make([]production.Entry, 0, len(records))
	for _, record := range records {
		entry, err := entryFromRecord(record)
		if err != nil {
			return nil, BatchError{RowNumber: record.RowNumber, Err: err}
		}
		entries = append(entries, entry)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entries...)

	return entries, nil
}

// Entries returns a snapshot of imported entries.
func (s *InMemoryStore) Entries() []production.Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]production.Entry(nil), s.entries...)
}
