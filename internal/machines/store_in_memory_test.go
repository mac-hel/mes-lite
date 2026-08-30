package machines

import (
	"errors"
	"fmt"
	"sync"
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

func TestInMemoryStoreRejectsDuplicateExternalEvent(t *testing.T) {
	store := NewInMemoryStore()
	event, err := NewEvent("machine-1", "external-1", EventTypeCycleCompleted, time.Now(), "sku-1", 2, "ws-1", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(t.Context(), event); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), event); !errors.Is(err, ErrDuplicateEvent) {
		t.Fatalf("expected ErrDuplicateEvent, got %v", err)
	}
}

func TestInMemoryStoreFindByExternalEventID(t *testing.T) {
	store := NewInMemoryStore()
	event, err := NewEvent("machine-1", "external-1", EventTypeCycleCompleted, time.Now(), "sku-1", 2, "ws-1", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), event); err != nil {
		t.Fatal(err)
	}

	found, err := store.FindByExternalEventID(t.Context(), "machine-1", "external-1")
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != event.ID {
		t.Fatalf("expected event %q, got %q", event.ID, found.ID)
	}

	_, err = store.FindByExternalEventID(t.Context(), "machine-1", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryStoreZeroValueSupportsConcurrentSaves(t *testing.T) {
	var store InMemoryStore
	const count = 100
	var wg sync.WaitGroup
	errs := make(chan error, count)

	for i := range count {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event, err := NewEvent("machine-1", fmt.Sprintf("external-%d", i), EventTypeCycleCompleted, time.Now(), "sku-1", 1, "ws-1", "")
			if err != nil {
				errs <- err
				return
			}
			errs <- store.Save(t.Context(), event)
		}(i)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	events, err := store.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != count {
		t.Fatalf("expected %d events, got %d", count, len(events))
	}
}
