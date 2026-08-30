package machines

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// Service coordinates machine event intake workflows.
type Service struct {
	store    Store
	counters intakeCounters
}

type intakeCounters struct {
	received         atomic.Uint64
	accepted         atomic.Uint64
	duplicateRetries atomic.Uint64
	conflicts        atomic.Uint64
	invalid          atomic.Uint64
}

// IntakeStats is a point-in-time snapshot of machine event intake counters.
type IntakeStats struct {
	Received         uint64
	Accepted         uint64
	DuplicateRetries uint64
	Conflicts        uint64
	Invalid          uint64
}

// NewService creates a machine integration service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// ReceiveEventCommand contains the event data received from a fake machine API call.
type ReceiveEventCommand struct {
	MachineID       string
	ExternalEventID string
	Type            EventType
	OccurredAt      time.Time
	ProductSKU      string
	Quantity        int
	Workstation     string
	Message         string
}

// ReceiveEvent validates and stores a machine event with idempotent retry handling.
func (s *Service) ReceiveEvent(ctx context.Context, cmd ReceiveEventCommand) (Event, error) {
	s.counters.received.Add(1)

	event, err := NewEvent(cmd.MachineID, cmd.ExternalEventID, cmd.Type, cmd.OccurredAt, cmd.ProductSKU, cmd.Quantity, cmd.Workstation, cmd.Message)
	if err != nil {
		s.counters.invalid.Add(1)
		return Event{}, err
	}

	if err := s.store.Save(ctx, event); err == nil {
		s.counters.accepted.Add(1)
		return event, nil
	} else if !errors.Is(err, ErrDuplicateEvent) {
		return Event{}, err
	}

	existing, err := s.store.FindByExternalEventID(ctx, event.MachineID, event.ExternalEventID)
	if err != nil {
		return Event{}, err
	}
	if !existing.samePayload(event) {
		s.counters.conflicts.Add(1)
		return Event{}, fmt.Errorf("machine event %q/%q has different payload: %w", event.MachineID, event.ExternalEventID, ErrEventConflict)
	}
	s.counters.duplicateRetries.Add(1)
	return existing, nil
}

// Stats returns a race-safe snapshot of machine event intake counters.
func (s *Service) Stats() IntakeStats {
	return IntakeStats{
		Received:         s.counters.received.Load(),
		Accepted:         s.counters.accepted.Load(),
		DuplicateRetries: s.counters.duplicateRetries.Load(),
		Conflicts:        s.counters.conflicts.Load(),
		Invalid:          s.counters.invalid.Load(),
	}
}
