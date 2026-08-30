package machines

import (
	"context"
	"sync"
)

// NewInMemoryStore creates an in-memory machine event store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

// InMemoryStore stores machine events in process memory.
type InMemoryStore struct {
	mu     sync.RWMutex
	events []Event
}

// Save stores one machine event.
func (s *InMemoryStore) Save(_ context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)
	return nil
}

// List returns a snapshot of received machine events.
func (s *InMemoryStore) List(_ context.Context) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]Event(nil), s.events...), nil
}
