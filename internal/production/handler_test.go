package production

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-fuego/fuego"

	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/products"
)

func testRegistrationService(t *testing.T, store Store) *Service {
	t.Helper()

	empStore := employees.NewInMemoryStore()
	prodStore := products.NewInMemoryStore()

	emp, err := employees.NewEmployee("emp-1", "Ana", "Worker", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := empStore.Save(t.Context(), emp); err != nil {
		t.Fatal(err)
	}
	prod, err := products.NewProduct("sku-1", "Ventilation Unit", "piece", products.CategoryVentilation)
	if err != nil {
		t.Fatal(err)
	}
	if err := prodStore.Save(t.Context(), prod); err != nil {
		t.Fatal(err)
	}

	return NewService(store, empStore, prodStore)
}

func TestHandler_Register(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(testRegistrationService(t, store))

	s := fuego.NewServer()
	fuego.Post(s, "/production-entries", handler.Register)

	body := `{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z","comment":"batch finished"}`
	req := httptest.NewRequest(http.MethodPost, "/production-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var entry Entry
	if err := json.NewDecoder(w.Body).Decode(&entry); err != nil {
		t.Fatal(err)
	}

	if entry.ID == "" {
		t.Fatal("expected generated ID")
	}
	if entry.EmployeeID != "emp-1" {
		t.Errorf("expected employee emp-1, got %q", entry.EmployeeID)
	}
	if entry.ProductSKU != "sku-1" {
		t.Errorf("expected product sku-1, got %q", entry.ProductSKU)
	}
	if entry.Quantity != 12 {
		t.Errorf("expected quantity 12, got %d", entry.Quantity)
	}
	if entry.Workstation != "ws-1" {
		t.Errorf("expected workstation ws-1, got %q", entry.Workstation)
	}
	if entry.Comment != "batch finished" {
		t.Errorf("expected comment, got %q", entry.Comment)
	}
	if !entry.Timestamp.Equal(time.Date(2026, 8, 8, 10, 30, 0, 0, time.UTC)) {
		t.Errorf("expected UTC timestamp, got %s", entry.Timestamp)
	}

	stored, err := store.FindByID(t.Context(), entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored != entry {
		t.Errorf("expected stored entry %#v, got %#v", entry, stored)
	}
}

func TestHandler_Register_Validation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"missing employee", `{"productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`},
		{"missing product", `{"employeeId":"emp-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`},
		{"zero quantity", `{"employeeId":"emp-1","productSku":"sku-1","quantity":0,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`},
		{"missing workstation", `{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"timestamp":"2026-08-08T10:30:00Z"}`},
		{"blank workstation", `{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":" ","timestamp":"2026-08-08T10:30:00Z"}`},
		{"missing timestamp", `{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1"}`},
	}

	store := NewInMemoryStore()
	handler := NewHandler(testRegistrationService(t, store))
	s := fuego.NewServer()
	fuego.Post(s, "/production-entries", handler.Register)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/production-entries", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_Register_MissingReference(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(testRegistrationService(t, store))
	s := fuego.NewServer()
	fuego.Post(s, "/production-entries", handler.Register)

	body := `{"employeeId":"missing","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/production-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_List(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(testRegistrationService(t, store))
	s := fuego.NewServer()
	fuego.Get(s, "/production-entries", handler.List)

	first := mustProductionEntry(t, "00000000-0000-4000-8000-000000000001", "emp-1", "sku-1", 12, "assembly-1", "2026-08-08T10:30:00Z")
	second := mustProductionEntry(t, "00000000-0000-4000-8000-000000000002", "emp-1", "sku-1", 7, "assembly-2", "2026-08-09T10:30:00Z")
	third := mustProductionEntry(t, "00000000-0000-4000-8000-000000000003", "emp-1", "sku-2", 4, "paint", "2026-08-10T10:30:00Z")
	for _, entry := range []Entry{first, second, third} {
		if err := store.Save(t.Context(), entry); err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/production-entries?productSku=sku-1&workstation=assembly&from=2026-08-08T00:00:00Z&to=2026-08-10T00:00:00Z&limit=1", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response ListProductionEntriesResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.Pagination.Count != 1 {
		t.Fatalf("expected one entry, got %d", response.Pagination.Count)
	}
	if response.Entries[0].ID != second.ID {
		t.Fatalf("expected newest matching entry %q, got %q", second.ID, response.Entries[0].ID)
	}
}

func TestHandler_List_InvalidQuery(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(testRegistrationService(t, store))
	s := fuego.NewServer()
	fuego.Get(s, "/production-entries", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/production-entries?from=not-a-time", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func mustProductionEntry(t *testing.T, id, employeeID, productSKU string, quantity int, workstation, timestamp string) Entry {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatal(err)
	}
	entry, err := NewEntry(id, employeeID, productSKU, quantity, workstation, parsed, "")
	if err != nil {
		t.Fatal(err)
	}
	return entry
}
