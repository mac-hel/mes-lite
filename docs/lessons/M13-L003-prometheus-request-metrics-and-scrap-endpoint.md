### Lesson 13.3 Scope

Expose Prometheus-compatible HTTP metrics and record low-cardinality request counters and durations.

#### Business Context

Operators need numeric signals that can be graphed and alerted on. Logs answer what happened for individual requests, but metrics answer aggregate questions such as how many requests are served and whether latency is increasing.

#### Problem

The application had structured request logs and correlation IDs, but no metrics endpoint. Prometheus could not scrape request volume or latency, and there was no foundation for future business or runtime metrics.

#### Design Discussion

Added `internal/platform/metrics` because metrics are technical observability plumbing, not a business slice. The package owns a Prometheus registry, HTTP request counter and request-duration histogram.

The first labels are intentionally low-cardinality: `method` and `status`. Route/path labels are postponed because raw paths can include IDs and create unbounded time series. If route-template labels are added later, they should use stable templates such as `/production-orders/{id}`, not raw URL paths.

`/metrics` is public and hidden from OpenAPI. In a real deployment, network policy or ingress rules would usually restrict it to the monitoring system. The endpoint is registered before request middleware, so metrics scrapes do not count as application traffic.

#### Go Concepts

- middleware composition around `http.Handler`
- response-writer wrapping reused for status capture
- histogram versus counter metric types
- dependency introduction when the standard library lacks the protocol tooling
- low-cardinality label design

#### Architecture Concepts

- observability platform package
- Prometheus pull model through `/metrics`
- metrics middleware separated from business handlers
- avoiding high-cardinality labels at the platform boundary

### Lesson 13.3 Completion Notes

#### Business Context

MES Lite now exposes Prometheus-compatible metrics for HTTP request volume and latency.

#### Problem

Operators could read request logs, but they could not scrape aggregate request metrics for dashboards or alerts.

#### Design Discussion

Introduced Prometheus through `github.com/prometheus/client_golang`. This dependency is justified because Prometheus exposition, registries and metric types are ecosystem standards and are not provided by the Go standard library.

The implementation uses a package-local registry instead of the global default registry. This keeps tests isolated and prevents duplicate metric registration panics when multiple servers are constructed in the same process.

HTTP metrics use only `method` and `status` labels. This avoids the common production mistake of labeling by raw path, where IDs such as `/production-orders/123` and `/production-orders/456` create separate time series.

#### Implementation

- Added `internal/platform/metrics`.
- Added `HTTPMetrics` with a Prometheus registry.
- Added `mes_lite_http_requests_total` counter.
- Added `mes_lite_http_request_duration_seconds` histogram.
- Added metrics middleware for request count and duration.
- Added `/metrics` endpoint using `promhttp.HandlerFor`.
- Registered `/metrics` before application middleware so scrapes do not count themselves.
- Hid `/metrics` from generated OpenAPI documentation.
- Added Prometheus client dependency.

#### Tests

- Added metrics middleware test for explicit status-code recording.
- Added metrics middleware test for default `200 OK` recording when handlers only write a body.
- Added server test proving `/metrics` returns request metrics after a `/health` request.
- Added server test coverage that `/metrics` scrape does not increment the application request counter.
- Verified with `go test ./internal/platform/metrics ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No business handlers changed. Metrics were added as platform middleware, keeping observability concerns out of the vertical slices.

#### Code Review

An experienced Go engineer would approve the small metrics foundation and isolated registry. The main positive review point is label restraint: method/status metrics are less detailed than route metrics, but they avoid unsafe cardinality until stable route-template labels are designed.

The main follow-up is tracing. Metrics can show that latency increased; traces will help explain where a slow request spent time.

#### Exercises

- Add a test proving a `404 Not Found` response increments the correct status label if routed through the same middleware path.
- Design safe route-template labels for Fuego routes without using raw URL paths.
- Decide whether `/metrics` should remain public or require network-level protection in deployment.

#### Interview Questions

- What is the difference between logs and metrics?
- Why are high-cardinality Prometheus labels dangerous?
- Why use a histogram for request duration instead of only logging duration?
- Why can a package-local Prometheus registry be better than the global default registry in tests?

#### Roadmap Update

- Lesson 13.3 completed.
- Current lesson moved to Lesson 13.4.
