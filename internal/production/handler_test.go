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
