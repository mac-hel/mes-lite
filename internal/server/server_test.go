package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/csvimport"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/machines"
	"github.com/mac-hel/mes-lite/internal/orders"
	"github.com/mac-hel/mes-lite/internal/platform/config"
	"github.com/mac-hel/mes-lite/internal/platform/jobs"
	platformlogging "github.com/mac-hel/mes-lite/internal/platform/logging"
	"github.com/mac-hel/mes-lite/internal/platform/version"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/reporting"
)

func testHandlers(t *testing.T) (*auth.Handler, *auth.Middleware, *auth.TokenManager, *employees.Handler, *products.Handler, *production.Handler, *orders.Handler, *reporting.Handler) {
	t.Helper()

	authStore := auth.NewInMemoryStore()
	user, err := auth.NewUser("user-1", "admin@example.com", "secret", auth.RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := authStore.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}
	tokens, err := auth.NewTokenManager("test-secret-with-at-least-32-characters")
	if err != nil {
		t.Fatal(err)
	}

	empStore := employees.NewInMemoryStore()
	prodStore := products.NewInMemoryStore()
	productionStore := production.NewInMemoryStore()

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

	productionService := production.NewService(productionStore, productionStore, empStore, prodStore)
	ordersStore := orders.NewInMemoryStore()
	ordersService := orders.NewService(ordersStore, empStore, prodStore)
	reportingStore := reporting.NewInMemoryStoreWithReports(
		[]reporting.DailyProductionRow{{
			Day:           mustTime(t, "2026-08-18T00:00:00Z"),
			ProductSKU:    "sku-1",
			TotalQuantity: 12,
			EntryCount:    1,
		}},
		[]reporting.EmployeeProductivityRow{{
			EmployeeID:    "emp-1",
			FirstName:     "Ana",
			LastName:      "Worker",
			TotalQuantity: 12,
			EntryCount:    1,
		}},
		[]reporting.ProductStatisticsRow{{
			ProductSKU:    "sku-1",
			ProductName:   "Ventilation Unit",
			TotalQuantity: 12,
			EntryCount:    1,
			EmployeeCount: 1,
		}},
		[]reporting.DailyEmployeeProductionRow{{
			Day:           mustTime(t, "2026-08-18T00:00:00Z"),
			ProductSKU:    "sku-1",
			ProductName:   "Ventilation Unit",
			EmployeeID:    "emp-1",
			FirstName:     "Ana",
			LastName:      "Worker",
			TotalQuantity: 12,
			EntryCount:    1,
		}},
		[]reporting.EmployeeProductivityProductRow{{
			EmployeeID:    "emp-1",
			FirstName:     "Ana",
			LastName:      "Worker",
			ProductSKU:    "sku-1",
			ProductName:   "Ventilation Unit",
			TotalQuantity: 12,
			EntryCount:    1,
		}},
	)

	return auth.NewHandler(auth.NewService(authStore, tokens)), auth.NewMiddleware(tokens), tokens, employees.NewHandler(empStore), products.NewHandler(prodStore), production.NewHandler(productionService, productionService), orders.NewHandler(ordersService), reporting.NewHandler(reportingStore)
}

func TestHealthEndpoint(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
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

func TestReadinessEndpointReady(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	var called bool
	s.SetReadinessCheck(func(ctx context.Context) error {
		called = true
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if !called {
		t.Fatal("expected readiness check to be called")
	}

	var body readyResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if body.Status != "ready" {
		t.Fatalf("expected status ready, got %q", body.Status)
	}
}

func TestReadinessEndpointNotReady(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	s.SetReadinessCheck(func(ctx context.Context) error {
		return errors.New("database unavailable")
	})

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", resp.StatusCode)
	}
}

func TestRequestLoggingAddsCorrelationID(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	var logs bytes.Buffer
	logger, err := platformlogging.New(&logs, "info", "json")
	if err != nil {
		t.Fatal(err)
	}
	s := NewWithLogger(config.Config{}, logger, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-ID", "server-request-1")
	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != "server-request-1" {
		t.Fatalf("expected response request ID header, got %q", w.Header().Get("X-Request-ID"))
	}

	var record map[string]any
	if err := json.Unmarshal(logs.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record["msg"] != "http request" {
		t.Fatalf("expected http request log, got %#v", record["msg"])
	}
	if record["request_id"] != "server-request-1" {
		t.Fatalf("expected request_id server-request-1, got %#v", record["request_id"])
	}
	if record["path"] != "/health" {
		t.Fatalf("expected path /health, got %#v", record["path"])
	}
	if record["status"] != float64(http.StatusOK) {
		t.Fatalf("expected status 200, got %#v", record["status"])
	}
}

func TestMetricsEndpoint(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)

	healthReq := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthW := httptest.NewRecorder()
	s.Mux.ServeHTTP(healthW, healthReq)
	if healthW.Code != http.StatusOK {
		t.Fatalf("expected health status 200, got %d", healthW.Code)
	}

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsW := httptest.NewRecorder()
	s.Mux.ServeHTTP(metricsW, metricsReq)
	if metricsW.Code != http.StatusOK {
		t.Fatalf("expected metrics status 200, got %d", metricsW.Code)
	}
	out := metricsW.Body.String()
	if !strings.Contains(out, `mes_lite_http_requests_total{method="GET",status="200"} 1`) {
		t.Fatalf("expected health request metric, got:\n%s", out)
	}
	if strings.Contains(out, `method="GET",status="200"} 2`) {
		t.Fatalf("expected metrics scrape not to count itself, got:\n%s", out)
	}
}

func TestVersionEndpoint(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	resp := w.Result()
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body version.Info
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body.Version != "dev" {
		t.Errorf("expected version 'dev', got %q", body.Version)
	}
	if body.Commit != "none" {
		t.Errorf("expected commit 'none', got %q", body.Commit)
	}
	if body.BuildTime != "unknown" {
		t.Errorf("expected build time 'unknown', got %q", body.BuildTime)
	}
}

func TestNotFound(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
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
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"requestId":"server-production-1","employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleAdmin)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntriesRouteRequiresAuthentication(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"requestId":"server-production-unauthenticated","employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLoginRouteRemainsPublic(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"email":"admin@example.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeCreateRequiresAdminRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"id":"emp-2","firstName":"Bob","lastName":"Worker","email":"bob@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductListAllowsLeaderRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntriesRouteAllowsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"requestId":"server-production-worker","employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntriesReviewAllowsLeaderRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"requestId":"server-production-review-seed","employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	setAuthorization(t, createReq, tokens, auth.RoleWorker)
	createW := httptest.NewRecorder()
	s.Mux.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusOK {
		t.Fatalf("expected create status 200, got %d: %s", createW.Code, createW.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/production-entries?employeeId=emp-1", nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntriesReviewForbidsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/production-entries", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntriesRouteIsIdempotentByRequestID(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"requestId":"server-production-idempotent","employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)

	firstReq := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	firstReq.Header.Set("Content-Type", "application/json")
	setAuthorization(t, firstReq, tokens, auth.RoleWorker)
	firstW := httptest.NewRecorder()
	s.Mux.ServeHTTP(firstW, firstReq)
	if firstW.Code != http.StatusOK {
		t.Fatalf("expected first status 200, got %d: %s", firstW.Code, firstW.Body.String())
	}
	var first production.Entry
	if err := json.NewDecoder(firstW.Body).Decode(&first); err != nil {
		t.Fatal(err)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	secondReq.Header.Set("Content-Type", "application/json")
	setAuthorization(t, secondReq, tokens, auth.RoleWorker)
	secondW := httptest.NewRecorder()
	s.Mux.ServeHTTP(secondW, secondReq)
	if secondW.Code != http.StatusOK {
		t.Fatalf("expected retry status 200, got %d: %s", secondW.Code, secondW.Body.String())
	}
	var second production.Entry
	if err := json.NewDecoder(secondW.Body).Decode(&second); err != nil {
		t.Fatal(err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected retry entry id %q, got %q", first.ID, second.ID)
	}
}

func TestProductionEntryCorrectionsRouteAllowsLeaderRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	entryID := createProductionEntry(t, s, tokens)
	body := []byte(`{"reason":"quantity typo","employeeId":"emp-1","productSku":"sku-1","quantity":13,"workstation":"ws-2","timestamp":"2026-08-08T11:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries/"+entryID+"/corrections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var correction production.Correction
	if err := json.NewDecoder(w.Body).Decode(&correction); err != nil {
		t.Fatal(err)
	}
	if correction.ActorUserID != "user-1" {
		t.Fatalf("expected actor user-1, got %q", correction.ActorUserID)
	}
}

func TestProductionEntryCorrectionsRouteForbidsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	entryID := createProductionEntry(t, s, tokens)
	body := []byte(`{"reason":"quantity typo","employeeId":"emp-1","productSku":"sku-1","quantity":13,"workstation":"ws-2","timestamp":"2026-08-08T11:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries/"+entryID+"/corrections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntryCorrectionsRouteListsHistory(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	entryID := createProductionEntry(t, s, tokens)
	body := []byte(`{"reason":"quantity typo","employeeId":"emp-1","productSku":"sku-1","quantity":13,"workstation":"ws-2","timestamp":"2026-08-08T11:30:00Z"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/production-entries/"+entryID+"/corrections", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	setAuthorization(t, createReq, tokens, auth.RoleManager)
	createW := httptest.NewRecorder()
	s.Mux.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusOK {
		t.Fatalf("expected correction status 200, got %d: %s", createW.Code, createW.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/production-entries/"+entryID+"/corrections", nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionOrdersRouteAllowsManagerCreate(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"lines":[{"productSku":"sku-1","plannedQuantity":12}]}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionOrdersRouteRejectsWorkerCreate(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"lines":[{"productSku":"sku-1","plannedQuantity":12}]}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionOrdersRouteAllowsLeaderRead(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"lines":[{"productSku":"sku-1","plannedQuantity":12}]}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleManager)
	createW := httptest.NewRecorder()
	s.Mux.ServeHTTP(createW, req)
	if createW.Code != http.StatusOK {
		t.Fatalf("expected create status 200, got %d: %s", createW.Code, createW.Body.String())
	}
	var createdOrder struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(createW.Body).Decode(&createdOrder); err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(http.MethodGet, "/production-orders/"+createdOrder.ID, nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionOrdersRouteRequiresAuthentication(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	body := []byte(`{"lines":[{"productSku":"sku-1","plannedQuantity":12}]}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionOrdersAssignmentRequiresManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	orderID := createProductionOrder(t, s, tokens)
	body := []byte(`{"employeeId":"emp-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders/"+orderID+"/assignments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionOrdersReleaseAllowsLeaderRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	orderID := createProductionOrder(t, s, tokens)
	body := []byte(`{"employeeId":"emp-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders/"+orderID+"/assignments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleManager)
	s.Mux.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodPut, "/production-orders/"+orderID+"/release", nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDailyProductionReportAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/daily-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDailyProductionReportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/daily-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeProductivityReportAllowsLeaderRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeProductivityReportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductStatisticsReportAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/product-statistics?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductStatisticsReportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/product-statistics?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDailyEmployeeProductionReportAllowsLeaderRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/daily-employee-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleLeader)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDailyEmployeeProductionReportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/daily-employee-production?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeProductivityProductsReportAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity/products?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeProductivityProductsReportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	req := httptest.NewRequest(http.MethodGet, "/reports/employee-productivity/products?from=2026-08-18T00:00:00Z&to=2026-08-19T00:00:00Z", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntryImportAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	csvImportH := csvimport.NewHandler(csvimport.NewService(csvimport.NewInMemoryStore()))
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH, csvImportH)
	body := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
	}, "\n")
	req := httptest.NewRequest(http.MethodPost, "/imports/production-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/csv")
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var summary csvimport.ImportSummary
	if err := json.NewDecoder(w.Body).Decode(&summary); err != nil {
		t.Fatal(err)
	}
	if summary.TotalRows != 1 || summary.ValidRows != 1 || summary.InvalidRows != 0 || summary.ImportedRows != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestProductionEntryImportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	csvImportH := csvimport.NewHandler(csvimport.NewService(csvimport.NewInMemoryStore()))
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH, csvImportH)
	body := "employee_id,product_sku,quantity,workstation,timestamp,comment\n"
	req := httptest.NewRequest(http.MethodPost, "/imports/production-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/csv")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntryImportRequiresAuthentication(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	csvImportH := csvimport.NewHandler(csvimport.NewService(csvimport.NewInMemoryStore()))
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH, csvImportH)
	body := "employee_id,product_sku,quantity,workstation,timestamp,comment\n"
	req := httptest.NewRequest(http.MethodPost, "/imports/production-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/csv")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductionEntryImportRejectsInvalidCSVHeader(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	csvImportH := csvimport.NewHandler(csvimport.NewService(csvimport.NewInMemoryStore()))
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH, csvImportH)
	req := httptest.NewRequest(http.MethodPost, "/imports/production-entries", strings.NewReader("employee_id,quantity\n"))
	req.Header.Set("Content-Type", "text/csv")
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAsyncProductionEntryImportAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	jobQueue := jobs.NewQueue(4)
	csvImportService := csvimport.NewService(csvimport.NewInMemoryStore())
	csvImportH := csvimport.NewHandlerWithAsync(csvImportService, csvimport.NewAsyncService(jobQueue, t.TempDir()))
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH, csvImportH)
	body := strings.Join([]string{
		"employee_id,product_sku,quantity,workstation,timestamp,comment",
		"emp-1,sku-1,12,ws-1,2026-08-20T10:00:00Z,valid",
	}, "\n")
	req := httptest.NewRequest(http.MethodPost, "/imports/production-entries/jobs", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/csv")
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var response jobs.JobResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.ID == "" || response.Type != jobs.TypeProductionEntryImport.String() || response.Status != jobs.StatusQueued.String() {
		t.Fatalf("unexpected job response: %+v", response)
	}
}

func TestAsyncProductionEntryImportRejectsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	jobQueue := jobs.NewQueue(4)
	csvImportService := csvimport.NewService(csvimport.NewInMemoryStore())
	csvImportH := csvimport.NewHandlerWithAsync(csvImportService, csvimport.NewAsyncService(jobQueue, t.TempDir()))
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH, csvImportH)
	req := httptest.NewRequest(http.MethodPost, "/imports/production-entries/jobs", strings.NewReader("employee_id,product_sku,quantity,workstation,timestamp,comment\n"))
	req.Header.Set("Content-Type", "text/csv")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestJobStatusRouteAllowsManagers(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	jobQueue, jobWorkers, jobsH := testJobHandlers(t)
	RegisterJobRoutes(s, authM, jobsH)

	if err := jobQueue.Enqueue(t.Context(), newServerTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/jobs/job-1", nil)
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var response jobs.JobResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.ID != "job-1" || response.Status != jobs.StatusQueued.String() {
		t.Fatalf("unexpected response: %+v", response)
	}

	if err := stopServerTestWorkers(t, jobWorkers); err != nil {
		t.Fatal(err)
	}
}

func TestJobStatusRouteRejectsWorkers(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	_, jobWorkers, jobsH := testJobHandlers(t)
	RegisterJobRoutes(s, authM, jobsH)

	req := httptest.NewRequest(http.MethodGet, "/jobs/job-1", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}

	if err := stopServerTestWorkers(t, jobWorkers); err != nil {
		t.Fatal(err)
	}
}

func TestCancelJobRouteAllowsAdmins(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	jobQueue, jobWorkers, jobsH := testJobHandlers(t)
	RegisterJobRoutes(s, authM, jobsH)

	if err := jobQueue.Enqueue(t.Context(), newServerTestJob(t, "job-1")); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/jobs/job-1/cancel", nil)
	setAuthorization(t, req, tokens, auth.RoleAdmin)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var response jobs.JobResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.Status != jobs.StatusCancelled.String() || !response.CancelRequested {
		t.Fatalf("expected cancelled response, got %+v", response)
	}

	if err := stopServerTestWorkers(t, jobWorkers); err != nil {
		t.Fatal(err)
	}
}

func TestMachineEventRouteAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	machineStore := machines.NewInMemoryStore()
	RegisterMachineRoutes(s, authM, machines.NewHandler(machines.NewService(machineStore)))

	body := []byte(`{"externalEventId":"machine-event-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":2,"workstation":"ws-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	events, err := machineStore.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 machine event, got %d", len(events))
	}
}

func TestMachineEventRouteForbidsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	machineStore := machines.NewInMemoryStore()
	RegisterMachineRoutes(s, authM, machines.NewHandler(machines.NewService(machineStore)))

	body := []byte(`{"externalEventId":"machine-event-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":2,"workstation":"ws-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMachineEventRouteRequiresAuthentication(t *testing.T) {
	authH, authM, _, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	machineStore := machines.NewInMemoryStore()
	RegisterMachineRoutes(s, authM, machines.NewHandler(machines.NewService(machineStore)))

	body := []byte(`{"externalEventId":"machine-event-1","type":"cycle_completed","occurredAt":"2026-08-30T10:30:00Z","productSku":"sku-1","quantity":2,"workstation":"ws-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/machines/machine-1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMachineStatsRouteAllowsManagerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	machineStore := machines.NewInMemoryStore()
	RegisterMachineRoutes(s, authM, machines.NewHandler(machines.NewService(machineStore)))

	req := httptest.NewRequest(http.MethodGet, "/machines/events/stats", nil)
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMachineStatsRouteForbidsWorkerRole(t *testing.T) {
	authH, authM, tokens, empH, prodH, productionH, ordersH, reportingH := testHandlers(t)
	s := New(config.Config{}, authH, authM, empH, prodH, productionH, ordersH, reportingH)
	machineStore := machines.NewInMemoryStore()
	RegisterMachineRoutes(s, authM, machines.NewHandler(machines.NewService(machineStore)))

	req := httptest.NewRequest(http.MethodGet, "/machines/events/stats", nil)
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func setAuthorization(t *testing.T, req *http.Request, tokens *auth.TokenManager, role auth.Role) {
	t.Helper()
	user, err := auth.NewUser("user-1", "user@example.com", "secret", role)
	if err != nil {
		t.Fatal(err)
	}
	token, err := tokens.Issue(user)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
}

func testJobHandlers(t *testing.T) (*jobs.Queue, *jobs.WorkerPool, *jobs.HTTPHandler) {
	t.Helper()

	queue := jobs.NewQueue(4)
	workers, err := jobs.NewWorkerPool(queue, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	return queue, workers, jobs.NewHTTPHandler(queue, workers)
}

func newServerTestJob(t *testing.T, id string) jobs.Job {
	t.Helper()

	job, err := jobs.NewJob(id, jobs.TypeProductionEntryImport, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return job
}

func stopServerTestWorkers(t *testing.T, workers *jobs.WorkerPool) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	return workers.Stop(ctx)
}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func createProductionEntry(t *testing.T, s *Server, tokens *auth.TokenManager) string {
	t.Helper()
	body := []byte(`{"requestId":"production-correction-seed","employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected create status 200, got %d: %s", w.Code, w.Body.String())
	}
	var created production.EntryResponse
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	return created.ID
}

func createProductionOrder(t *testing.T, s *Server, tokens *auth.TokenManager) string {
	t.Helper()
	body := []byte(`{"lines":[{"productSku":"sku-1","plannedQuantity":12}]}`)
	req := httptest.NewRequest(http.MethodPost, "/production-orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleManager)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected create status 200, got %d: %s", w.Code, w.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	return created.ID
}
