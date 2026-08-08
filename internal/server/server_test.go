package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mac-hel/mes-lite/internal/config"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/products"
)

func testHandlers(t *testing.T) (*employees.Handler, *products.Handler, *production.Handler) {
	t.Helper()

	empStore := employees.NewInMemoryStore()
	prodStore := products.NewInMemoryStore()
	productionStore := production.NewInMemoryStore()

	if err := empStore.Save(t.Context(), employees.NewEmployee("emp-1", "Ana", "Worker", "ana@example.com")); err != nil {
		t.Fatal(err)
	}
	if err := prodStore.Save(t.Context(), products.NewProduct("sku-1", "Ventilation Unit", "piece", products.CategoryVentilation)); err != nil {
		t.Fatal(err)
	}

	productionService := production.NewService(productionStore, empStore, prodStore)

	return employees.NewHandler(empStore), products.NewHandler(prodStore), production.NewHandler(productionService)
}

func TestHealthEndpoint(t *testing.T) {
	empH, prodH, productionH := testHandlers(t)
	s := New(config.Config{}, empH, prodH, productionH)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", body.Status)
	}
}

func TestVersionEndpoint(t *testing.T) {
	empH, prodH, productionH := testHandlers(t)
	s := New(config.Config{}, empH, prodH, productionH)
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body versionResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body.Version != "dev" {
		t.Errorf("expected version 'dev', got %q", body.Version)
	}
}

func TestNotFound(t *testing.T) {
	empH, prodH, productionH := testHandlers(t)
	s := New(config.Config{}, empH, prodH, productionH)
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestProductionEntriesRoute(t *testing.T) {
	empH, prodH, productionH := testHandlers(t)
	s := New(config.Config{}, empH, prodH, productionH)
	body := []byte(`{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}
