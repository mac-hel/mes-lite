package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPMetricsRecordsRequests(t *testing.T) {
	metrics := NewHTTPMetrics()
	handler := metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/employees", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsW := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsW, metricsReq)
	out := metricsW.Body.String()

	if !strings.Contains(out, `mes_lite_http_requests_total{method="POST",status="201"} 1`) {
		t.Fatalf("expected request counter, got:\n%s", out)
	}
	if !strings.Contains(out, `mes_lite_http_request_duration_seconds_bucket{method="POST",status="201"`) {
		t.Fatalf("expected duration histogram, got:\n%s", out)
	}
}

func TestHTTPMetricsDefaultStatusOK(t *testing.T) {
	metrics := NewHTTPMetrics()
	handler := metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsW := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsW, metricsReq)
	out := metricsW.Body.String()

	if !strings.Contains(out, `mes_lite_http_requests_total{method="GET",status="200"} 1`) {
		t.Fatalf("expected default 200 counter, got:\n%s", out)
	}
}
