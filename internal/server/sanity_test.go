package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/config"
)

func TestSanity_MVPAPIFlow(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)

	loginBody := []byte(`{"email":"admin@example.com","password":"secret"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	s.Mux.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d: %s", loginW.Code, loginW.Body.String())
	}

	var loginResult auth.LoginResult
	if err := json.NewDecoder(loginW.Body).Decode(&loginResult); err != nil {
		t.Fatal(err)
	}
	if loginResult.AccessToken == "" {
		t.Fatal("expected login to return an access token")
	}

	productionBody := []byte(`{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	productionReq := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(productionBody))
	productionReq.Header.Set("Content-Type", "application/json")
	productionReq.Header.Set("Authorization", "Bearer "+loginResult.AccessToken)
	productionW := httptest.NewRecorder()
	s.Mux.ServeHTTP(productionW, productionReq)
	if productionW.Code != http.StatusOK {
		t.Fatalf("expected production registration status 200, got %d: %s", productionW.Code, productionW.Body.String())
	}

	reportReq := httptest.NewRequest(http.MethodGet, "/reports/daily-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	reportReq.Header.Set("Authorization", "Bearer "+loginResult.AccessToken)
	reportW := httptest.NewRecorder()
	s.Mux.ServeHTTP(reportW, reportReq)
	if reportW.Code != http.StatusOK {
		t.Fatalf("expected daily production report status 200, got %d: %s", reportW.Code, reportW.Body.String())
	}
}
