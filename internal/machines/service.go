package machines

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Service coordinates machine event intake workflows.
type Service struct {
	store Store
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
	event, err := NewEvent(cmd.MachineID, cmd.ExternalEventID, cmd.Type, cmd.OccurredAt, cmd.ProductSKU, cmd.Quantity, cmd.Workstation, cmd.Message)
	if err != nil {
		return Event{}, err
	}

	if err := s.store.Save(ctx, event); err == nil {
		return event, nil
	} else if !errors.Is(err, ErrDuplicateEvent) {
		return Event{}, err
	}

	existing, err := s.store.FindByExternalEventID(ctx, event.MachineID, event.ExternalEventID)
	if err != nil {
		return Event{}, err
	}
	if !existing.samePayload(event) {
		return Event{}, fmt.Errorf("machine event %q/%q has different payload: %w", event.MachineID, event.ExternalEventID, ErrEventConflict)
	}
	return existing, nil
}
