package machines

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mac-hel/mes-lite/internal/platform/ids"
)

// ErrInvalidEvent marks invalid machine event input.
var ErrInvalidEvent = errors.New("invalid machine event")

// EventType identifies the kind of signal received from a production machine.
type EventType string

const (
	// EventTypeCycleCompleted means a machine reports completed production output.
	EventTypeCycleCompleted EventType = "cycle_completed"
	// EventTypeStateChanged means a machine reports an operational state change.
	EventTypeStateChanged EventType = "state_changed"
	// EventTypeAlarmRaised means a machine reports an alarm or fault condition.
	EventTypeAlarmRaised EventType = "alarm_raised"
)

// Valid reports whether the event type is supported by MES Lite.
func (t EventType) Valid() bool {
	switch t {
	case EventTypeCycleCompleted, EventTypeStateChanged, EventTypeAlarmRaised:
		return true
	default:
		return false
	}
}

// Event is one signal received from a production machine.
type Event struct {
	ID              string
	MachineID       string
	ExternalEventID string
	Type            EventType
	OccurredAt      time.Time
	ProductSKU      string
	Quantity        int
	Workstation     string
	Message         string
}

// NewEvent creates a machine event and normalizes text fields and timestamps.
func NewEvent(machineID, externalEventID string, eventType EventType, occurredAt time.Time, productSKU string, quantity int, workstation, message string) (Event, error) {
	event := Event{
		ID:              ids.New(),
		MachineID:       strings.TrimSpace(machineID),
		ExternalEventID: strings.TrimSpace(externalEventID),
		Type:            eventType,
		OccurredAt:      occurredAt.UTC(),
		ProductSKU:      strings.TrimSpace(productSKU),
		Quantity:        quantity,
		Workstation:     strings.TrimSpace(workstation),
		Message:         strings.TrimSpace(message),
	}
	if err := event.Validate(); err != nil {
		return Event{}, err
	}
	return event, nil
}

// Validate checks the invariants common to all machine events.
func (e Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.MachineID) == "" {
		return fmt.Errorf("machine id is required: %w", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.ExternalEventID) == "" {
		return fmt.Errorf("external event id is required: %w", ErrInvalidEvent)
	}
	if !e.Type.Valid() {
		return fmt.Errorf("event type %q is not supported: %w", e.Type, ErrInvalidEvent)
	}
	if e.OccurredAt.IsZero() {
		return fmt.Errorf("occurred at is required: %w", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.Workstation) == "" {
		return fmt.Errorf("workstation is required: %w", ErrInvalidEvent)
	}
	if e.Type == EventTypeCycleCompleted {
		if strings.TrimSpace(e.ProductSKU) == "" {
			return fmt.Errorf("product sku is required for cycle completed events: %w", ErrInvalidEvent)
		}
		if e.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive for cycle completed events: %w", ErrInvalidEvent)
		}
	}
	return nil
}
