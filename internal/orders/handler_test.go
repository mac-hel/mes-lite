package orders

import (
	"context"
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

func TestHandler_Create(t *testing.T) {
	handler, _ := testHandler(t)
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

	handler, _ := testHandler(t)
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
	handler, _ := testHandler(t)
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

	duplicateHandler := NewHandler(alwaysDuplicateService{})
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

	handler := NewHandler(NewService(store, seededEmployeeStore(t), seededProductStore(t)))
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
	handler, _ := testHandler(t)
	s := fuego.NewServer()
	fuego.Get(s, "/production-orders/{id}", handler.Get)
	req := httptest.NewRequest(http.MethodGet, "/production-orders/missing", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_AssignEmployee(t *testing.T) {
	handler, store := testHandler(t)
	order := mustPersistedOrder(t, store)
	s := fuego.NewServer()
	fuego.Post(s, "/production-orders/{id}/assignments", handler.AssignEmployee)

	req := httptest.NewRequest(http.MethodPost, "/production-orders/"+order.ID()+"/assignments", strings.NewReader(`{"employeeId":"emp-1"}`))
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
	if len(resp.AssignedEmployeeIDs) != 1 || resp.AssignedEmployeeIDs[0] != "emp-1" {
		t.Fatalf("assigned employees = %#v, want [emp-1]", resp.AssignedEmployeeIDs)
	}
}

func TestHandler_ReleaseStartComplete(t *testing.T) {
	handler, store := testHandler(t)
	order := mustPersistedOrder(t, store)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := store.Update(t.Context(), order); err != nil {
		t.Fatal(err)
	}
	s := fuego.NewServer()
	fuego.Put(s, "/production-orders/{id}/release", handler.Release)
	fuego.Put(s, "/production-orders/{id}/start", handler.Start)
	fuego.Put(s, "/production-orders/{id}/complete", handler.Complete)

	assertTransition(t, s, "/production-orders/"+order.ID()+"/release", StatusReleased)
	assertTransition(t, s, "/production-orders/"+order.ID()+"/start", StatusInProgress)
	assertTransition(t, s, "/production-orders/"+order.ID()+"/complete", StatusCompleted)
}

func TestHandler_Cancel(t *testing.T) {
	handler, store := testHandler(t)
	order := mustPersistedOrder(t, store)
	s := fuego.NewServer()
	fuego.Put(s, "/production-orders/{id}/cancel", handler.Cancel)

	assertTransition(t, s, "/production-orders/"+order.ID()+"/cancel", StatusCancelled)
}

func testNow() time.Time {
	return time.Date(2026, 8, 18, 10, 0, 0, 0, time.UTC)
}

func mustPersistedOrder(t *testing.T, store *InMemoryStore) Order {
	t.Helper()
	order, err := NewOrder("order-1", mustOrderLines(t, mustOrderLine(t, "shaft-1", 2)), testNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), order); err != nil {
		t.Fatal(err)
	}
	return order
}

func assertTransition(t *testing.T, s *fuego.Server, path string, want Status) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, path, nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp OrderResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != want {
		t.Fatalf("Status = %q, want %q", resp.Status, want)
	}
}

func isUUIDShaped(id string) bool {
	return len(id) == 36 && id[8] == '-' && id[13] == '-' && id[18] == '-' && id[23] == '-'
}

func testHandler(t *testing.T) (*Handler, *InMemoryStore) {
	t.Helper()
	store := NewInMemoryStore()
	return NewHandler(NewService(store, seededEmployeeStore(t), seededProductStore(t))), store
}

func seededEmployeeStore(t *testing.T) *employees.InMemoryStore {
	t.Helper()
	store := employees.NewInMemoryStore()
	emp, err := employees.NewEmployee("emp-1", "Ana", "Worker", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), emp); err != nil {
		t.Fatal(err)
	}
	return store
}

func seededProductStore(t *testing.T) *products.InMemoryStore {
	t.Helper()
	store := products.NewInMemoryStore()
	for _, product := range []struct {
		sku      string
		name     string
		category products.ProductCategory
	}{
		{"shaft-1", "Shaft", products.CategoryVentilation},
		{"filter-1", "Filter", products.CategoryFilter},
	} {
		prod, err := products.NewProduct(product.sku, product.name, "piece", product.category)
		if err != nil {
			t.Fatal(err)
		}
		if err := store.Save(t.Context(), prod); err != nil {
			t.Fatal(err)
		}
	}
	return store
}

type alwaysDuplicateService struct{}

func (alwaysDuplicateService) Create(ctx context.Context, cmd CreateCommand) (Order, error) {
	return Order{}, ErrAlreadyExists
}

func (alwaysDuplicateService) Get(ctx context.Context, id string) (Order, error) {
	return Order{}, ErrNotFound
}

func (alwaysDuplicateService) AssignEmployee(ctx context.Context, cmd AssignEmployeeCommand) (Order, error) {
	return Order{}, ErrNotFound
}

func (alwaysDuplicateService) Release(ctx context.Context, id string) (Order, error) {
	return Order{}, ErrNotFound
}

func (alwaysDuplicateService) Start(ctx context.Context, id string) (Order, error) {
	return Order{}, ErrNotFound
}

func (alwaysDuplicateService) Complete(ctx context.Context, id string) (Order, error) {
	return Order{}, ErrNotFound
}

func (alwaysDuplicateService) Cancel(ctx context.Context, id string) (Order, error) {
	return Order{}, ErrNotFound
}
