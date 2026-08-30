package machines

import (
	"testing"
	"time"
)

func TestInMemoryStoreSaveAndList(t *testing.T) {
	store := NewInMemoryStore()
	event, err := NewEvent("machine-1", "external-1", EventTypeCycleCompleted, time.Now(), "sku-1", 2, "ws-1", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(t.Context(), event); err != nil {
		t.Fatal(err)
	}

	events, err := store.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ID != event.ID {
		t.Fatalf("expected event %q, got %q", event.ID, events[0].ID)
	}

	events[0].MachineID = "mutated"
	again, err := store.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if again[0].MachineID != "machine-1" {
		t.Fatalf("expected store snapshot to be isolated, got %q", again[0].MachineID)
	}
}
