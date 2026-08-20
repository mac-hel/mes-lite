package reporting

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-fuego/fuego"
)

func TestHandler_DailyProduction(t *testing.T) {
	store := NewInMemoryStore([]DailyProductionRow{{
		Day:           time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
		ProductSKU:    "shaft-1",
		TotalQuantity: 7,
		EntryCount:    2,
	}})
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/reports/daily-production", handler.DailyProduction)

	req := httptest.NewRequest(http.MethodGet, "/reports/daily-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp DailyProductionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Rows) != 1 {
		t.Fatalf("row count = %d, want 1", len(resp.Rows))
	}
	row := resp.Rows[0]
	if row.ProductSKU != "shaft-1" || row.TotalQuantity != 7 || row.EntryCount != 2 {
		t.Fatalf("row = %#v, want shaft-1 quantity 7 count 2", row)
	}
}

func TestHandler_DailyProductionInvalidRange(t *testing.T) {
	handler := NewHandler(NewInMemoryStore(nil))
	s := fuego.NewServer()
	fuego.Get(s, "/reports/daily-production", handler.DailyProduction)

	tests := []struct {
		name string
		url  string
	}{
		{name: "missing from", url: "/reports/daily-production?to=2026-08-19T00:00:00Z"},
		{name: "invalid from", url: "/reports/daily-production?from=2026-08-18&to=2026-08-19T00:00:00Z"},
		{name: "from equals to", url: "/reports/daily-production?from=2026-08-18T00:00:00Z&to=2026-08-18T00:00:00Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_EmployeeProductivity(t *testing.T) {
	store := NewInMemoryStoreWithReports(nil, []EmployeeProductivityRow{{
		EmployeeID:    "emp-1",
		FirstName:     "Ana",
		LastName:      "Worker",
		TotalQuantity: 7,
		EntryCount:    2,
	}}, nil, nil, nil)
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/reports/employee-productivity", handler.EmployeeProductivity)

	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp EmployeeProductivityResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Rows) != 1 {
		t.Fatalf("row count = %d, want 1", len(resp.Rows))
	}
	row := resp.Rows[0]
	if row.EmployeeID != "emp-1" || row.FirstName != "Ana" || row.TotalQuantity != 7 || row.EntryCount != 2 {
		t.Fatalf("row = %#v, want employee emp-1 quantity 7 count 2", row)
	}
}

func TestHandler_EmployeeProductivityInvalidRange(t *testing.T) {
	handler := NewHandler(NewInMemoryStoreWithReports(nil, nil, nil, nil, nil))
	s := fuego.NewServer()
	fuego.Get(s, "/reports/employee-productivity", handler.EmployeeProductivity)

	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity?from=2026-08-18T00:00:00Z&to=2026-08-18T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ProductStatistics(t *testing.T) {
	store := NewInMemoryStoreWithReports(nil, nil, []ProductStatisticsRow{{
		ProductSKU:    "shaft-1",
		ProductName:   "Shaft",
		TotalQuantity: 7,
		EntryCount:    2,
		EmployeeCount: 2,
	}}, nil, nil)
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/reports/product-statistics", handler.ProductStatistics)

	req := httptest.NewRequest(http.MethodGet, "/reports/product-statistics?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp ProductStatisticsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Rows) != 1 {
		t.Fatalf("row count = %d, want 1", len(resp.Rows))
	}
	row := resp.Rows[0]
	if row.ProductSKU != "shaft-1" || row.ProductName != "Shaft" || row.TotalQuantity != 7 || row.EmployeeCount != 2 {
		t.Fatalf("row = %#v, want shaft-1 quantity 7 employee count 2", row)
	}
}

func TestHandler_ProductStatisticsInvalidRange(t *testing.T) {
	handler := NewHandler(NewInMemoryStoreWithReports(nil, nil, nil, nil, nil))
	s := fuego.NewServer()
	fuego.Get(s, "/reports/product-statistics", handler.ProductStatistics)

	req := httptest.NewRequest(http.MethodGet, "/reports/product-statistics?from=2026-08-18T00:00:00Z&to=2026-08-18T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DailyEmployeeProduction(t *testing.T) {
	store := NewInMemoryStoreWithReports(nil, nil, nil, []DailyEmployeeProductionRow{{
		Day:           time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
		ProductSKU:    "shaft-1",
		ProductName:   "Shaft",
		EmployeeID:    "emp-1",
		FirstName:     "Ana",
		LastName:      "Worker",
		TotalQuantity: 7,
		EntryCount:    2,
	}}, nil)
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/reports/daily-employee-production", handler.DailyEmployeeProduction)

	req := httptest.NewRequest(http.MethodGet, "/reports/daily-employee-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp DailyEmployeeProductionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Rows) != 1 {
		t.Fatalf("row count = %d, want 1", len(resp.Rows))
	}
	row := resp.Rows[0]
	if row.ProductSKU != "shaft-1" || row.EmployeeID != "emp-1" || row.TotalQuantity != 7 || row.EntryCount != 2 {
		t.Fatalf("row = %#v, want shaft-1 employee emp-1 quantity 7 count 2", row)
	}
}

func TestHandler_EmployeeProductivityProducts(t *testing.T) {
	store := NewInMemoryStoreWithReports(nil, nil, nil, nil, []EmployeeProductivityProductRow{{
		EmployeeID:    "emp-1",
		FirstName:     "Ana",
		LastName:      "Worker",
		ProductSKU:    "shaft-1",
		ProductName:   "Shaft",
		TotalQuantity: 7,
		EntryCount:    2,
	}})
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Get(s, "/reports/employee-productivity/products", handler.EmployeeProductivityProducts)

	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity/products?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp EmployeeProductivityProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Rows) != 1 {
		t.Fatalf("row count = %d, want 1", len(resp.Rows))
	}
	row := resp.Rows[0]
	if row.EmployeeID != "emp-1" || row.ProductSKU != "shaft-1" || row.TotalQuantity != 7 || row.EntryCount != 2 {
		t.Fatalf("row = %#v, want employee emp-1 shaft-1 quantity 7 count 2", row)
	}
}

func TestHandler_DetailedReportsInvalidRange(t *testing.T) {
	handler := NewHandler(NewInMemoryStoreWithReports(nil, nil, nil, nil, nil))
	s := fuego.NewServer()
	fuego.Get(s, "/reports/daily-employee-production", handler.DailyEmployeeProduction)
	fuego.Get(s, "/reports/employee-productivity/products", handler.EmployeeProductivityProducts)

	tests := []string{
		"/reports/daily-employee-production?from=2026-08-18T00:00:00Z&to=2026-08-18T00:00:00Z",
		"/reports/employee-productivity/products?from=2026-08-18T00:00:00Z&to=2026-08-18T00:00:00Z",
	}
	for _, url := range tests {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected status 400, got %d: %s", url, w.Code, w.Body.String())
		}
	}
}
