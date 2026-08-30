package machines

import (
	"errors"
	"testing"
	"time"
)

func TestServiceReceiveEventStoresNewEvent(t *testing.T) {
	store := NewInMemoryStore()
	service := NewService(store)

	event, err := service.ReceiveEvent(t.Context(), validReceiveEventCommand())
	if err != nil {
		t.Fatal(err)
	}
	if event.ID == "" {
		t.Fatal("expected generated event id")
	}

	events, err := store.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 stored event, got %d", len(events))
	}
}

func TestServiceReceiveEventReturnsExistingEventForIdenticalRetry(t *testing.T) {
	store := NewInMemoryStore()
	service := NewService(store)
	cmd := validReceiveEventCommand()

	first, err := service.ReceiveEvent(t.Context(), cmd)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.ReceiveEvent(t.Context(), cmd)
	if err != nil {
		t.Fatal(err)
	}

	if second.ID != first.ID {
		t.Fatalf("expected retry to return original event ID %q, got %q", first.ID, second.ID)
	}
	events, err := store.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one stored event after retry, got %d", len(events))
	}
}

func TestServiceReceiveEventRejectsConflictingRetry(t *testing.T) {
	service := NewService(NewInMemoryStore())
	cmd := validReceiveEventCommand()

	if _, err := service.ReceiveEvent(t.Context(), cmd); err != nil {
		t.Fatal(err)
	}
	cmd.Quantity = 99

	_, err := service.ReceiveEvent(t.Context(), cmd)
	if !errors.Is(err, ErrEventConflict) {
		t.Fatalf("expected ErrEventConflict, got %v", err)
	}
}

func TestServiceReceiveEventRejectsInvalidEvent(t *testing.T) {
	service := NewService(NewInMemoryStore())
	cmd := validReceiveEventCommand()
	cmd.MachineID = ""

	_, err := service.ReceiveEvent(t.Context(), cmd)
	if !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("expected ErrInvalidEvent, got %v", err)
	}
}

func validReceiveEventCommand() ReceiveEventCommand {
	return ReceiveEventCommand{
		MachineID:       "machine-1",
		ExternalEventID: "external-1",
		Type:            EventTypeCycleCompleted,
		OccurredAt:      time.Date(2026, 8, 30, 10, 30, 0, 0, time.UTC),
		ProductSKU:      "sku-1",
		Quantity:        3,
		Workstation:     "ws-1",
		Message:         "ok",
	}
}
