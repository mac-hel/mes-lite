package machines

import (
	"errors"
	"testing"
	"time"
)

func TestNewEventNormalizesFields(t *testing.T) {
	occurredAt := time.Date(2026, 8, 30, 12, 30, 0, 0, time.FixedZone("UTC+2", 2*60*60))

	event, err := NewEvent(" machine-1 ", " external-1 ", EventTypeCycleCompleted, occurredAt, " sku-1 ", 12, " ws-1 ", " completed ")
	if err != nil {
		t.Fatal(err)
	}

	if event.ID == "" {
		t.Fatal("expected generated event ID")
	}
	if event.MachineID != "machine-1" {
		t.Fatalf("expected trimmed machine id, got %q", event.MachineID)
	}
	if event.ExternalEventID != "external-1" {
		t.Fatalf("expected trimmed external event id, got %q", event.ExternalEventID)
	}
	if event.ProductSKU != "sku-1" {
		t.Fatalf("expected trimmed product sku, got %q", event.ProductSKU)
	}
	if event.Workstation != "ws-1" {
		t.Fatalf("expected trimmed workstation, got %q", event.Workstation)
	}
	if event.Message != "completed" {
		t.Fatalf("expected trimmed message, got %q", event.Message)
	}
	if event.OccurredAt.Location() != time.UTC {
		t.Fatalf("expected UTC timestamp, got %s", event.OccurredAt.Location())
	}
}

func TestNewEventValidation(t *testing.T) {
	tests := []struct {
		name     string
		makeArgs func() (string, string, EventType, time.Time, string, int, string, string)
	}{
		{
			name: "missing machine id",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "", "external-1", EventTypeCycleCompleted, time.Now(), "sku-1", 1, "ws-1", ""
			},
		},
		{
			name: "missing external event id",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "machine-1", "", EventTypeCycleCompleted, time.Now(), "sku-1", 1, "ws-1", ""
			},
		},
		{
			name: "unsupported type",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "machine-1", "external-1", EventType("unknown"), time.Now(), "sku-1", 1, "ws-1", ""
			},
		},
		{
			name: "missing timestamp",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "machine-1", "external-1", EventTypeCycleCompleted, time.Time{}, "sku-1", 1, "ws-1", ""
			},
		},
		{
			name: "missing workstation",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "machine-1", "external-1", EventTypeCycleCompleted, time.Now(), "sku-1", 1, "", ""
			},
		},
		{
			name: "missing cycle product",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "machine-1", "external-1", EventTypeCycleCompleted, time.Now(), "", 1, "ws-1", ""
			},
		},
		{
			name: "non-positive cycle quantity",
			makeArgs: func() (string, string, EventType, time.Time, string, int, string, string) {
				return "machine-1", "external-1", EventTypeCycleCompleted, time.Now(), "sku-1", 0, "ws-1", ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machineID, externalEventID, eventType, occurredAt, productSKU, quantity, workstation, message := tt.makeArgs()
			_, err := NewEvent(machineID, externalEventID, eventType, occurredAt, productSKU, quantity, workstation, message)
			if !errors.Is(err, ErrInvalidEvent) {
				t.Fatalf("expected ErrInvalidEvent, got %v", err)
			}
		})
	}
}

func TestStateChangedEventDoesNotRequireProductionFields(t *testing.T) {
	_, err := NewEvent("machine-1", "external-1", EventTypeStateChanged, time.Now(), "", 0, "ws-1", "idle")
	if err != nil {
		t.Fatal(err)
	}
}
