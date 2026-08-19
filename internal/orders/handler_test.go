package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-fuego/fuego"
)

func TestHandler_Create(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Post(s, "/production-orders", handler.Create)

	body := `{"lines":[{"productSku":"shaft-1","plannedQuantity":2},{"productSku":"filter-1","plannedQuantity":4}],"assignedEmployeeIds":["emp-1"]}`
	req := httptest.NewRequest(http.MethodPost, "/production-orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp OrderResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if !isUUIDShaped(resp.ID) {
		t.Fatalf("ID = %q, want UUID-shaped generated ID", resp.ID)
	}
	if len(resp.Lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(resp.Lines))
	}
	if resp.Lines[0].ProductSKU != "shaft-1" || resp.Lines[0].PlannedQuantity != 2 {
		t.Fatalf("first line = %#v, want shaft-1 quantity 2", resp.Lines[0])
	}
	if len(resp.AssignedEmployeeIDs) != 1 || resp.AssignedEmployeeIDs[0] != "emp-1" {
		t.Fatalf("assigned employees = %#v, want [emp-1]", resp.AssignedEmployeeIDs)
	}
}

func TestHandler_CreateValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"missing lines", `{}`},
		{"blank product sku", `{"lines":[{"productSku":" ","plannedQuantity":2}]}`},
		{"zero quantity", `{"lines":[{"productSku":"shaft-1","plannedQuantity":0}]}`},
		{"duplicate product", `{"lines":[{"productSku":"shaft-1","plannedQuantity":2},{"productSku":"shaft-1","plannedQuantity":4}]}`},
		{"blank assigned employee", `{"lines":[{"productSku":"shaft-1","plannedQuantity":2}],"assignedEmployeeIds":[" "]}`},
	}

	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Post(s, "/production-orders", handler.Create)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/production-orders", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_CreateDuplicate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Post(s, "/production-orders", handler.Create)

	body := `{"lines":[{"productSku":"shaft-1","plannedQuantity":2}]}`
	req := httptest.NewRequest(http.MethodPost, "/production-orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected first create status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp OrderResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	order, err := store.FindByID(t.Context(), resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), order); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("Save() error = %v, want ErrAlreadyExists", err)
	}

	duplicateHandler := NewHandler(alwaysDuplicateStore{})
	s = fuego.NewServer()
	fuego.Post(s, "/production-orders", duplicateHandler.Create)
	req = httptest.NewRequest(http.MethodPost, "/production-orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Get(t *testing.T) {
	store := NewInMemoryStore()
	line := mustOrderLine(t, "shaft-1", 2)
	order, err := NewOrder("order-1", mustOrderLines(t, line), testNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), order); err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/production-orders/{id}", handler.Get)
	req := httptest.NewRequest(http.MethodGet, "/production-orders/order-1", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp OrderResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != "order-1" {
		t.Fatalf("ID = %q, want order-1", resp.ID)
	}
}

func TestHandler_GetNotFound(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/production-orders/{id}", handler.Get)
	req := httptest.NewRequest(http.MethodGet, "/production-orders/missing", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func testNow() time.Time {
	return time.Date(2026, 8, 18, 10, 0, 0, 0, time.UTC)
}

func isUUIDShaped(id string) bool {
	return len(id) == 36 && id[8] == '-' && id[13] == '-' && id[18] == '-' && id[23] == '-'
}

type alwaysDuplicateStore struct{}

func (alwaysDuplicateStore) Save(ctx context.Context, order Order) error {
	return ErrAlreadyExists
}

func (alwaysDuplicateStore) FindByID(ctx context.Context, id string) (Order, error) {
	return Order{}, ErrNotFound
}
