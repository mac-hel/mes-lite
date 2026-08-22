package production

import (
	"context"
	"fmt"
	"sort"
	"strings"
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

// List returns production entries matching review filters, newest first.
func (s *InMemoryStore) List(_ context.Context, opts ListOptions) ([]Entry, error) {
	opts, err := opts.normalize()
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(s.entries))
	for _, entry := range s.entries {
		if opts.EmployeeID != "" && entry.EmployeeID != opts.EmployeeID {
			continue
		}
		if opts.ProductSKU != "" && entry.ProductSKU != opts.ProductSKU {
			continue
		}
		if opts.Workstation != "" && !strings.Contains(strings.ToLower(entry.Workstation), strings.ToLower(opts.Workstation)) {
			continue
		}
		if !opts.From.IsZero() && entry.Timestamp.Before(opts.From) {
			continue
		}
		if !opts.To.IsZero() && !entry.Timestamp.Before(opts.To) {
			continue
		}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Timestamp.Equal(entries[j].Timestamp) {
			return entries[i].ID > entries[j].ID
		}
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	if opts.Offset >= len(entries) {
		return []Entry{}, nil
	}
	end := opts.Offset + opts.Limit
	if end > len(entries) {
		end = len(entries)
	}

	return entries[opts.Offset:end], nil
}
