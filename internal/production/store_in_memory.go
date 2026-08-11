package production

import (
	"context"
	"fmt"
)

// NewInMemoryStore creates a map-based in-memory [InMemoryStore].
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{entries: make(map[string]Entry)}
}

// InMemoryStore is a map-based in-memory implementation of [Store].
type InMemoryStore struct {
	entries map[string]Entry
}

// Save stores a production entry keyed by ID.
func (s *InMemoryStore) Save(_ context.Context, entry Entry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	if _, ok := s.entries[entry.ID]; ok {
		return fmt.Errorf("production entry %q: %w", entry.ID, ErrAlreadyExists)
	}
	s.entries[entry.ID] = entry
	return nil
}

// FindByID looks up a production entry by ID. Returns [ErrNotFound] if not found.
func (s *InMemoryStore) FindByID(_ context.Context, id string) (Entry, error) {
	entry, ok := s.entries[id]
	if !ok {
		return Entry{}, fmt.Errorf("production entry %q: %w", id, ErrNotFound)
	}
	return entry, nil
}
