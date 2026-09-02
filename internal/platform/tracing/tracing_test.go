package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestNewProviderRejectsInvalidExporter(t *testing.T) {
	if _, err := NewProvider(t.Context(), Config{Exporter: "zipkin"}); err == nil {
		t.Fatal("expected invalid exporter error")
	}
}

func TestMiddlewareCreatesServerSpan(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	handler := Middleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := TraceIDFromContext(r.Context())
		if !ok {
			t.Fatal("expected trace ID in request context")
		}
		if traceID == "" {
			t.Fatal("expected non-empty trace ID")
		}
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/production-entries?ignored=true", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected one ended span, got %d", len(spans))
	}
	span := spans[0]
	if span.Name() != "POST /production-entries" {
		t.Fatalf("expected span name POST /production-entries, got %q", span.Name())
	}
	assertAttr(t, span.Attributes(), "http.request.method", http.MethodPost)
	assertAttr(t, span.Attributes(), "url.path", "/production-entries")
	assertAttr(t, span.Attributes(), "http.response.status_code", int64(http.StatusCreated))
}

func TestTraceIDFromContextWithoutSpan(t *testing.T) {
	if traceID, ok := TraceIDFromContext(context.Background()); ok || traceID != "" {
		t.Fatalf("expected no trace ID, got %q", traceID)
	}
}

func assertAttr(t *testing.T, attrs []attribute.KeyValue, key string, want any) {
	t.Helper()

	for _, attr := range attrs {
		if string(attr.Key) != key {
			continue
		}
		switch want := want.(type) {
		case string:
			if attr.Value.AsString() != want {
				t.Fatalf("expected %s=%q, got %q", key, want, attr.Value.AsString())
			}
		case int64:
			if attr.Value.AsInt64() != want {
				t.Fatalf("expected %s=%d, got %d", key, want, attr.Value.AsInt64())
			}
		default:
			t.Fatalf("unsupported attr assertion type %T", want)
		}
		return
	}
	t.Fatalf("expected attribute %s", key)
}
