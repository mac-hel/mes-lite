package machines

import (
	"context"
	"fmt"
	"sync"
)

// NewInMemoryStore creates an in-memory machine event store.
func NewInMemoryStore() Store {
	return &inMemoryStore{byExternalID: make(map[string]Event)}
}

type inMemoryStore struct {
	mu           sync.RWMutex
	events       []Event
	byExternalID map[string]Event
}

// Save stores one machine event.
func (s *inMemoryStore) Save(_ context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := externalEventKey(event.MachineID, event.ExternalEventID)
	if _, ok := s.byExternalID[key]; ok {
		return fmt.Errorf("machine event %q/%q: %w", event.MachineID, event.ExternalEventID, ErrDuplicateEvent)
	}

	s.events = append(s.events, event)
	s.byExternalID[key] = event
	return nil
}

// FindByExternalEventID looks up a machine event by machine and external event IDs.
func (s *inMemoryStore) FindByExternalEventID(_ context.Context, machineID, externalEventID string) (Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.byExternalID[externalEventKey(machineID, externalEventID)]
	if !ok {
		return Event{}, fmt.Errorf("machine event %q/%q: %w", machineID, externalEventID, ErrNotFound)
	}
	return event, nil
}

// List returns a snapshot of received machine events.
func (s *inMemoryStore) List(_ context.Context) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]Event(nil), s.events...), nil
}

func externalEventKey(machineID, externalEventID string) string {
	return machineID + "\x00" + externalEventID
}
