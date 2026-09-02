### Lesson 13.4 Scope

Add OpenTelemetry HTTP tracing so each request creates a server span and request logs can include a trace ID when tracing is active.

#### Business Context

Metrics can show that requests are slow or failing, but operators need traces to understand where a request spent time as the service grows. A trace ID also gives support and developers a common handle that can connect logs, client reports and distributed traces.

#### Problem

The application had structured logs, correlation IDs and Prometheus metrics, but no tracing foundation. Request logs could identify one request, but there was no OpenTelemetry span context that future database, job or integration spans could attach to.

#### Design Discussion

Tracing is added as platform middleware. Each HTTP request starts a server span with method, path, status and duration attributes. The middleware extracts W3C `traceparent` headers so upstream trace context can flow into MES Lite.

The provider supports `OTEL_TRACES_EXPORTER=none` by default and `OTEL_TRACES_EXPORTER=stdout` for local inspection. OTLP export is intentionally postponed. Adding a collector endpoint and deployment configuration is a production-integration concern, while this lesson focuses on request span creation and context propagation.

The middleware order is metrics, tracing, then request logging. This means request logs are emitted while the span context is still active, so logs can include `trace_id`.

#### Go Concepts

- OpenTelemetry tracer providers and spans
- context propagation through HTTP middleware
- W3C trace context extraction
- span attributes and status
- graceful tracer-provider shutdown

#### Architecture Concepts

- tracing as platform infrastructure
- request spans as the root for future child spans
- trace IDs connected to structured logs
- exporter choice kept in configuration

### Lesson 13.4 Completion Notes

#### Business Context

MES Lite now creates OpenTelemetry server spans for HTTP requests and can connect request logs to traces through `trace_id`.

#### Problem

Logs and metrics were useful but incomplete. There was no trace context for following one request through future internal operations.

#### Design Discussion

Added `internal/platform/tracing` with a configurable tracer provider and HTTP middleware. The middleware starts one server span per request and records low-risk request attributes. It avoids business-specific spans for now because database, job and integration tracing should be added where those boundaries are reviewed.

`cmd/server` configures the provider during startup and shuts it down with a timeout on exit. The default exporter is `none`, so local development does not emit trace payloads unexpectedly. `stdout` is available for learning and manual inspection.

Request logging now checks the active span context and includes `trace_id` when one exists.

#### Implementation

- Added `internal/platform/tracing`.
- Added `tracing.Config` and `NewProvider`.
- Added `OTEL_TRACES_EXPORTER` configuration with default `none`.
- Added `stdout` trace exporter support for local inspection.
- Installed W3C trace-context propagation.
- Added tracing HTTP middleware.
- Added span attributes for method, path, status and duration.
- Added trace IDs to request logs when an active span exists.
- Wired tracing startup and shutdown in `cmd/server`.

#### Tests

- Added invalid exporter configuration test.
- Added tracing middleware test using OpenTelemetry's span recorder.
- Added assertions for span name, method, path and response status attributes.
- Added test for missing trace ID when no span exists.
- Added request-log test proving `trace_id` is included when span context is active.
- Updated configuration tests for `OTEL_TRACES_EXPORTER`.
- Verified with `go test ./internal/platform/tracing ./internal/platform/logging ./internal/platform/config ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No business package was changed. Tracing was added as platform middleware and request logging gained only trace-context awareness.

#### Code Review

An experienced Go engineer would approve the foundation because tracing context is established at the HTTP boundary and exporter behavior is explicit. The main limitation is deliberate: OTLP export and child spans for database/job operations are not wired yet.

#### Exercises

- Send a request with an existing `traceparent` header and verify the server span joins the upstream trace.
- Enable `OTEL_TRACES_EXPORTER=stdout` and inspect the emitted span fields.
- Design where database query spans should be introduced without leaking sqlc details into handlers.

#### Interview Questions

- What does a trace show that logs and metrics do not?
- Why is context propagation required for distributed tracing?
- What is the difference between a trace ID and a span ID?
- Why should tracer-provider shutdown be part of graceful shutdown?

#### Roadmap Update

- Lesson 13.4 completed.
- Current lesson moved to Lesson 13.5.
- Known technical debt updated for missing OTLP exporter wiring.
