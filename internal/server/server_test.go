package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/config"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/orders"
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

	productionService := production.NewService(productionStore, empStore, prodStore)
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
	)

	return auth.NewHandler(auth.NewService(authStore, tokens)), auth.NewMiddleware(tokens), tokens, employees.NewHandler(empStore), products.NewHandler(prodStore), production.NewHandler(productionService), orders.NewHandler(ordersService), reporting.NewHandler(reportingStore)
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
	body := []byte(`{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
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
	body := []byte(`{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
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
	body := []byte(`{"employeeId":"emp-1","productSku":"sku-1","quantity":12,"workstation":"ws-1","timestamp":"2026-08-08T10:30:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/production-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthorization(t, req, tokens, auth.RoleWorker)
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
