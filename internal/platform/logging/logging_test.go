package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	logger, err := New(&buf, "info", "json")
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("database connected", "component", "server")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["msg"] != "database connected" {
		t.Fatalf("expected structured message, got %#v", got["msg"])
	}
	if got["component"] != "server" {
		t.Fatalf("expected structured component, got %#v", got["component"])
	}
}

func TestNewTextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger, err := New(&buf, "debug", "text")
	if err != nil {
		t.Fatal(err)
	}

	logger.Debug("debug enabled", "component", "test")

	out := buf.String()
	if !strings.Contains(out, "msg=\"debug enabled\"") {
		t.Fatalf("expected text slog output, got %q", out)
	}
	if !strings.Contains(out, "component=test") {
		t.Fatalf("expected structured component, got %q", out)
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name string
		want slog.Level
	}{
		{name: "", want: slog.LevelInfo},
		{name: "debug", want: slog.LevelDebug},
		{name: "info", want: slog.LevelInfo},
		{name: "warn", want: slog.LevelWarn},
		{name: "warning", want: slog.LevelWarn},
		{name: "error", want: slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevel(tt.name)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestInvalidConfiguration(t *testing.T) {
	if _, err := New(&bytes.Buffer{}, "verbose", "json"); err == nil {
		t.Fatal("expected invalid level error")
	}
	if _, err := New(&bytes.Buffer{}, "info", "xml"); err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestRequestLoggerUsesProvidedRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger, err := New(&buf, "info", "json")
	if err != nil {
		t.Fatal(err)
	}

	var gotRequestID string
	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := RequestIDFromContext(r.Context())
		if !ok {
			t.Fatal("expected request ID in context")
		}
		gotRequestID = requestID
		w.WriteHeader(http.StatusCreated)
	}))
	req := httptest.NewRequest(http.MethodPost, "/employees?ignored=true", nil)
	req.Header.Set(requestIDHeader, "request-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if gotRequestID != "request-123" {
		t.Fatalf("expected context request ID, got %q", gotRequestID)
	}
	if w.Header().Get(requestIDHeader) != "request-123" {
		t.Fatalf("expected response request ID header, got %q", w.Header().Get(requestIDHeader))
	}

	log := decodeLog(t, buf.Bytes())
	if log["request_id"] != "request-123" {
		t.Fatalf("expected request_id request-123, got %#v", log["request_id"])
	}
	if log["method"] != http.MethodPost {
		t.Fatalf("expected method POST, got %#v", log["method"])
	}
	if log["path"] != "/employees" {
		t.Fatalf("expected path without query string, got %#v", log["path"])
	}
	if log["status"] != float64(http.StatusCreated) {
		t.Fatalf("expected status 201, got %#v", log["status"])
	}
}

func TestRequestLoggerGeneratesRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger, err := New(&buf, "info", "json")
	if err != nil {
		t.Fatal(err)
	}

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := RequestIDFromContext(r.Context()); !ok {
			t.Fatal("expected generated request ID in context")
		}
	}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	requestID := w.Header().Get(requestIDHeader)
	if len(requestID) != 36 {
		t.Fatalf("expected generated UUID-shaped request ID, got %q", requestID)
	}

	log := decodeLog(t, buf.Bytes())
	if log["request_id"] != requestID {
		t.Fatalf("expected log request ID %q, got %#v", requestID, log["request_id"])
	}
}

func decodeLog(t *testing.T, data []byte) map[string]any {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	return got
}
