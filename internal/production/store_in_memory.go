package production

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// NewInMemoryStore creates a map-based in-memory [InMemoryStore].
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		entries:     make(map[string]Entry),
		corrections: make(map[string]Correction),
	}
}

// InMemoryStore is a map-based in-memory implementation of [Store].
type InMemoryStore struct {
	entries     map[string]Entry
	corrections map[string]Correction
}

// Save stores a production entry keyed by ID.
func (s *InMemoryStore) Save(_ context.Context, entry Entry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	if _, ok := s.entries[entry.ID]; ok {
		return fmt.Errorf("production entry %q: %w", entry.ID, ErrAlreadyExists)
	}
	if entry.RequestID != "" {
		for _, existing := range s.entries {
			if existing.RequestID == entry.RequestID {
				return fmt.Errorf("production request id %q: %w", entry.RequestID, ErrRequestConflict)
			}
		}
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

// FindByRequestID looks up a production entry by idempotency request ID.
func (s *InMemoryStore) FindByRequestID(_ context.Context, requestID string) (Entry, error) {
	for _, entry := range s.entries {
		if entry.RequestID == requestID {
			return entry, nil
		}
	}
	return Entry{}, fmt.Errorf("production request id %q: %w", requestID, ErrNotFound)
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

// SaveCorrection stores an append-only correction keyed by correction ID.
func (s *InMemoryStore) SaveCorrection(_ context.Context, correction Correction) error {
	if err := correction.Validate(); err != nil {
		return err
	}
	if _, ok := s.entries[correction.EntryID]; !ok {
		return fmt.Errorf("production entry %q: %w", correction.EntryID, ErrNotFound)
	}
	if _, ok := s.corrections[correction.ID]; ok {
		return fmt.Errorf("production correction %q: %w", correction.ID, ErrAlreadyExists)
	}
	s.corrections[correction.ID] = correction
	return nil
}

// ListCorrections returns correction history for one production entry, newest first.
func (s *InMemoryStore) ListCorrections(_ context.Context, entryID string) ([]Correction, error) {
	if _, ok := s.entries[entryID]; !ok {
		return nil, fmt.Errorf("production entry %q: %w", entryID, ErrNotFound)
	}
	corrections := make([]Correction, 0)
	for _, correction := range s.corrections {
		if correction.EntryID == entryID {
			corrections = append(corrections, correction)
		}
	}
	sort.Slice(corrections, func(i, j int) bool {
		if corrections[i].CreatedAt.Equal(corrections[j].CreatedAt) {
			return corrections[i].ID > corrections[j].ID
		}
		return corrections[i].CreatedAt.After(corrections[j].CreatedAt)
	})
	return corrections, nil
}
