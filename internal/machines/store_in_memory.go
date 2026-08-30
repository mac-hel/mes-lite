package machines

import (
	"context"
	"fmt"
	"sync"
)

// NewInMemoryStore creates an in-memory machine event store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

// InMemoryStore stores machine events in process memory.
type InMemoryStore struct {
	initOnce     sync.Once
	mu           sync.RWMutex
	events       []Event
	byExternalID map[string]Event
}

func (s *InMemoryStore) init() {
	s.byExternalID = make(map[string]Event)
}

// Save stores one machine event.
func (s *InMemoryStore) Save(_ context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	s.initOnce.Do(s.init)

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
func (s *InMemoryStore) FindByExternalEventID(_ context.Context, machineID, externalEventID string) (Event, error) {
	s.initOnce.Do(s.init)

	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.byExternalID[externalEventKey(machineID, externalEventID)]
	if !ok {
		return Event{}, fmt.Errorf("machine event %q/%q: %w", machineID, externalEventID, ErrNotFound)
	}
	return event, nil
}

// List returns a snapshot of received machine events.
func (s *InMemoryStore) List(_ context.Context) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]Event(nil), s.events...), nil
}

func externalEventKey(machineID, externalEventID string) string {
	return machineID + "\x00" + externalEventID
}
