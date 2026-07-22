package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mac-hel/mes-lite/internal/config"
)

func TestHealthEndpoint(t *testing.T) {
	s := New(config.Config{})
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
	s := New(config.Config{})
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
	s := New(config.Config{})
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
