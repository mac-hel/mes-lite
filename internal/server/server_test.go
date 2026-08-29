package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/csvimport"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/orders"
	"github.com/mac-hel/mes-lite/internal/platform/config"
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

	var body versionResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body.Version != "dev" {
		t.Errorf("expected version 'dev', got %q", body.Version)
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
