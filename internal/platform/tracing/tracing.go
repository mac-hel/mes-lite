package tracing

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Config controls OpenTelemetry tracing setup.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Exporter       string
}

// NewProvider builds and installs the process OpenTelemetry tracer provider.
func NewProvider(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	serviceName := strings.TrimSpace(cfg.ServiceName)
	if serviceName == "" {
		serviceName = "mes-lite"
	}

	resource, err := sdkresource.Merge(
		sdkresource.Default(),
		sdkresource.NewWithAttributes("",
			attribute.String("service.name", serviceName),
			attribute.String("service.version", strings.TrimSpace(cfg.ServiceVersion)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create trace resource: %w", err)
	}

	options := []sdktrace.TracerProviderOption{sdktrace.WithResource(resource)}
	switch strings.ToLower(strings.TrimSpace(cfg.Exporter)) {
	case "", "none":
	case "stdout":
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("create stdout trace exporter: %w", err)
		}
		options = append(options, sdktrace.WithBatcher(exporter))
	default:
		return nil, fmt.Errorf("invalid trace exporter %q", cfg.Exporter)
	}

	provider := sdktrace.NewTracerProvider(options...)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return provider, nil
}

// Middleware starts one server span for each HTTP request.
func Middleware(provider oteltrace.TracerProvider) func(http.Handler) http.Handler {
	if provider == nil {
		provider = otel.GetTracerProvider()
	}
	tracer := provider.Tracer("mes-lite/http")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			spanName := r.Method + " " + r.URL.Path
			ctx, span := tracer.Start(ctx, spanName,
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				oteltrace.WithAttributes(
					attribute.String("http.request.method", r.Method),
					attribute.String("url.path", r.URL.Path),
				),
			)
			defer span.End()

			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			started := time.Now()
			next.ServeHTTP(recorder, r.WithContext(ctx))

			span.SetAttributes(
				attribute.Int("http.response.status_code", recorder.status),
				attribute.Int64("http.server.duration_ms", time.Since(started).Milliseconds()),
			)
			if recorder.status >= http.StatusInternalServerError {
				span.SetStatus(codes.Error, http.StatusText(recorder.status))
			}
		})
	}
}

// TraceIDFromContext returns the active span trace ID when one exists.
func TraceIDFromContext(ctx context.Context) (string, bool) {
	spanContext := oteltrace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return "", false
	}
	return spanContext.TraceID().String(), true
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(data)
}
