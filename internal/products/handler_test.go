package products

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-fuego/fuego"
)

func TestHandler_Create(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)

	body := `{"sku":"VX-100","name":"Ventilation Unit X100","unit":"piece","category":0}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	responseBody := w.Body.Bytes()
	var p ProductResponse
	if err := json.Unmarshal(responseBody, &p); err != nil {
		t.Fatal(err)
	}

	if p.SKU != "VX-100" {
		t.Errorf("expected SKU VX-100, got %q", p.SKU)
	}
	if p.Name != "Ventilation Unit X100" {
		t.Errorf("expected Name Ventilation Unit X100, got %q", p.Name)
	}
	if p.Unit != "piece" {
		t.Errorf("expected Unit piece, got %q", p.Unit)
	}
	if p.Category != CategoryVentilation {
		t.Errorf("expected Category 0 (Ventilation), got %d", p.Category)
	}
	if !p.IsActive {
		t.Error("expected IsActive true")
	}

	var fields map[string]any
	if err := json.Unmarshal(responseBody, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["sku"]; !ok {
		t.Fatal("expected lowercase sku field")
	}
	if _, ok := fields["SKU"]; ok {
		t.Fatal("did not expect capitalized SKU field")
	}
}

func TestHandler_DuplicateCreate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)

	body := `{"sku":"VX-100","name":"Test","unit":"piece","category":0}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	w := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Create_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing sku", `{"name":"Test","unit":"piece","category":0}`, http.StatusBadRequest},
		{"missing name", `{"sku":"001","unit":"piece","category":0}`, http.StatusBadRequest},
		{"missing unit", `{"sku":"001","name":"Test","category":0}`, http.StatusBadRequest},
		{"blank sku", `{"sku":" ","name":"Test","unit":"piece","category":0}`, http.StatusBadRequest},
		{"invalid category", `{"sku":"001","name":"Test","unit":"piece","category":99}`, http.StatusBadRequest},
	}

	store := NewInMemoryStore()
	handler := NewHandler(store)
	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_List_Empty(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Get(s, "/products", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Products == nil {
		t.Error("expected non-nil slice for empty list")
	}
	if len(resp.Products) != 0 {
		t.Errorf("expected 0 products, got %d", len(resp.Products))
	}
}

func TestHandler_List_AfterCreate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Get(s, "/products", handler.List)

	createProduct := func(sku, name string) {
		body := `{"sku":"` + sku + `","name":"` + name + `","unit":"piece","category":0}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d", sku, w.Code)
		}
	}

	createProduct("VX-100", "Ventilation X100")
	createProduct("FT-200", "Filter T200")
	createProduct("DC-300", "Duct C300")

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Products) != 3 {
		t.Fatalf("expected 3 products, got %d", len(resp.Products))
	}

	skus := make(map[string]bool, 3)
	for _, p := range resp.Products {
		skus[p.SKU] = true
	}
	for _, expected := range []string{"VX-100", "FT-200", "DC-300"} {
		if !skus[expected] {
			t.Errorf("expected product %s in list", expected)
		}
	}
}

func TestHandler_Update(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Put(s, "/products/{sku}", handler.Update)

	createBody := `{"sku":"VX-100","name":"Old Name","unit":"piece","category":0}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	updateBody := `{"name":"Updated Name","unit":"set","category":1,"version":1}`
	req = httptest.NewRequest(http.MethodPut, "/products/VX-100", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var p ProductResponse
	if err := json.NewDecoder(w.Body).Decode(&p); err != nil {
		t.Fatal(err)
	}

	if p.SKU != "VX-100" {
		t.Errorf("expected SKU VX-100, got %q", p.SKU)
	}
	if p.Name != "Updated Name" {
		t.Errorf("expected Name Updated Name, got %q", p.Name)
	}
	if p.Unit != "set" {
		t.Errorf("expected Unit set, got %q", p.Unit)
	}
	if p.Category != CategoryFilter {
		t.Errorf("expected Category Filter (1), got %d", p.Category)
	}
	if p.Version != 2 {
		t.Errorf("expected Version 2, got %d", p.Version)
	}
}

func TestHandler_Update_VersionConflict(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Put(s, "/products/{sku}", handler.Update)

	createBody := `{"sku":"VX-100","name":"Old Name","unit":"piece","category":0}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	updateBody := `{"name":"Updated Name","unit":"set","category":1,"version":1}`
	req = httptest.NewRequest(http.MethodPut, "/products/VX-100", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodPut, "/products/VX-100", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Update_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing name", `{"unit":"piece","category":0}`, http.StatusBadRequest},
		{"missing unit", `{"name":"Test","category":0}`, http.StatusBadRequest},
		{"blank name", `{"name":" ","unit":"piece","category":0}`, http.StatusBadRequest},
		{"invalid category", `{"name":"Test","unit":"piece","category":99}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryStore()
			handler := NewHandler(store)

			s := fuego.NewServer()
			fuego.Post(s, "/products", handler.Create)
			fuego.Put(s, "/products/{sku}", handler.Update)

			createBody := `{"sku":"VX-100","name":"Old Name","unit":"piece","category":0}`
			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(createBody))
			req.Header.Set("Content-Type", "application/json")
			s.Mux.ServeHTTP(httptest.NewRecorder(), req)

			req = httptest.NewRequest(http.MethodPut, "/products/VX-100", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_Update_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Put(s, "/products/{sku}", handler.Update)

	body := `{"name":"Test","unit":"piece","category":0}`
	req := httptest.NewRequest(http.MethodPut, "/products/NONEXISTENT", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Deactivate(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Put(s, "/products/{sku}/deactivate", handler.Deactivate)

	createBody := `{"sku":"VX-100","name":"Test","unit":"piece","category":0}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodPut, "/products/VX-100/deactivate", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var p ProductResponse
	if err := json.NewDecoder(w.Body).Decode(&p); err != nil {
		t.Fatal(err)
	}

	if p.SKU != "VX-100" {
		t.Errorf("expected SKU VX-100, got %q", p.SKU)
	}
	if p.IsActive {
		t.Error("expected IsActive false")
	}
}

func TestHandler_Deactivate_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Put(s, "/products/{sku}/deactivate", handler.Deactivate)

	req := httptest.NewRequest(http.MethodPut, "/products/NONEXISTENT/deactivate", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Search(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Get(s, "/products/search", handler.Search)

	createProduct := func(sku, name string, category ProductCategory) {
		categoryJSON := `"category":` + string(rune('0'+category))
		body := `{"sku":"` + sku + `","name":"` + name + `","unit":"piece",` + categoryJSON + `}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d: %s", sku, w.Code, w.Body.String())
		}
	}

	createProduct("VX-100", "Ventilation Unit X100", CategoryVentilation)
	createProduct("FT-200", "Filter Type 200", CategoryFilter)
	createProduct("DC-300", "Duct Connector 300", CategoryDuct)

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=vent", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Products) != 1 {
		t.Fatalf("expected 1 result for 'vent', got %d", len(resp.Products))
	}
	if resp.Products[0].SKU != "VX-100" {
		t.Errorf("expected VX-100, got %q", resp.Products[0].SKU)
	}
}

func TestHandler_Search_ByCategoryName(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Get(s, "/products/search", handler.Search)

	createProduct := func(sku, name string, category ProductCategory) {
		categoryJSON := `"category":` + string(rune('0'+category))
		body := `{"sku":"` + sku + `","name":"` + name + `","unit":"piece",` + categoryJSON + `}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d", sku, w.Code)
		}
	}

	createProduct("VX-100", "Unit One", CategoryVentilation)
	createProduct("VX-200", "Unit Two", CategoryVentilation)
	createProduct("FT-100", "Filter One", CategoryFilter)

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=filter", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Products) != 1 {
		t.Fatalf("expected 1 result for 'filter', got %d", len(resp.Products))
	}
	if resp.Products[0].SKU != "FT-100" {
		t.Errorf("expected FT-100, got %q", resp.Products[0].SKU)
	}
}

func TestHandler_Search_NoResults(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Get(s, "/products/search", handler.Search)

	createProduct := func(sku, name string, category ProductCategory) {
		categoryJSON := `"category":` + string(rune('0'+category))
		body := `{"sku":"` + sku + `","name":"` + name + `","unit":"piece",` + categoryJSON + `}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d", sku, w.Code)
		}
	}

	createProduct("VX-100", "Ventilation Unit", CategoryVentilation)

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=zzzzz", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Products == nil {
		t.Error("expected non-nil slice for empty results")
	}
	if len(resp.Products) != 0 {
		t.Errorf("expected 0 results, got %d", len(resp.Products))
	}
}

func TestHandler_Search_EmptyQuery(t *testing.T) {
	store := NewInMemoryStore()
	handler := NewHandler(store)

	s := fuego.NewServer()
	fuego.Post(s, "/products", handler.Create)
	fuego.Get(s, "/products/search", handler.Search)

	createProduct := func(sku, name string, category ProductCategory) {
		categoryJSON := `"category":` + string(rune('0'+category))
		body := `{"sku":"` + sku + `","name":"` + name + `","unit":"piece",` + categoryJSON + `}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("create %s: expected 200, got %d", sku, w.Code)
		}
	}

	createProduct("VX-100", "Unit", CategoryVentilation)
	createProduct("FT-200", "Filter", CategoryFilter)

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Products) != 2 {
		t.Errorf("expected all 2 products for empty query, got %d", len(resp.Products))
	}
}
