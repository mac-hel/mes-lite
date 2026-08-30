package machines

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego"
)

func TestCreateEvent(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(NewService(store))
	server := fuego.NewServer()
	fuego.Post(server, "/machines/{machineId}/events", handler.CreateEvent)

	body := []byte(`{"externalEventId":"external-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":3,"workstation":"ws-1","message":"ok"}`)
	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response EventResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.ID == "" {
		t.Fatal("expected generated event id")
	}
	if response.MachineID != "machine-1" {
		t.Fatalf("expected machine id from path, got %q", response.MachineID)
	}
	if response.ExternalEventID != "external-1" {
		t.Fatalf("expected external event id, got %q", response.ExternalEventID)
	}

	events, err := store.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 stored event, got %d", len(events))
	}
}

func TestCreateEventRejectsInvalidEvent(t *testing.T) {
	handler := NewHandler(NewService(NewInMemoryStore()))
	server := fuego.NewServer()
	fuego.Post(server, "/machines/{machineId}/events", handler.CreateEvent)

	body := []byte(`{"externalEventId":"external-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"","quantity":3,"workstation":"ws-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateEventReturnsExistingEventForIdenticalRetry(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(NewService(store))
	server := fuego.NewServer()
	fuego.Post(server, "/machines/{machineId}/events", handler.CreateEvent)
	body := []byte(`{"externalEventId":"external-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":3,"workstation":"ws-1"}`)

	first := postMachineEvent(t, server, body)
	second := postMachineEvent(t, server, body)

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

func TestCreateEventReturnsConflictForDifferentRetryPayload(t *testing.T) {
	handler := NewHandler(NewService(NewInMemoryStore()))
	server := fuego.NewServer()
	fuego.Post(server, "/machines/{machineId}/events", handler.CreateEvent)

	body := []byte(`{"externalEventId":"external-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":3,"workstation":"ws-1"}`)
	_ = postMachineEvent(t, server, body)

	conflictingBody := []byte(`{"externalEventId":"external-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":4,"workstation":"ws-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(conflictingBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func postMachineEvent(t *testing.T, server *fuego.Server, body []byte) EventResponse {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response EventResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	return response
}
