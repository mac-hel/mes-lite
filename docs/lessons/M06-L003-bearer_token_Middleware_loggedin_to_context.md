### Lesson 6.3 Scope

Add authentication middleware that verifies `Authorization: Bearer <token>` and stores the authenticated principal in request context.

#### Business Context

After login returns a signed token, business endpoints must reject anonymous callers. Workers, leaders and administrators should prove identity before touching production or master data.

#### Problem

JWTs could be issued and verified manually, but no request pipeline enforced token verification. Every business endpoint was still public.

#### Design Discussion

Authentication is implemented as route-scoped middleware using Fuego's `OptionMiddleware`. Public endpoints stay public: health, readiness, version and login. Business endpoints require a valid bearer token.

The middleware uses standard `net/http` shapes: `func(http.Handler) http.Handler`. This is idiomatic Go because middleware composes around the standard library instead of depending on framework-specific magic.

#### Go Concepts

- HTTP middleware as function composition
- request headers and `Authorization: Bearer`
- request context values for request-scoped identity
- unexported context key type to avoid collisions

#### Architecture Concepts

- authentication happens before handlers
- handlers remain focused on business work
- principal propagation is request-scoped, not global state
- role checks are postponed to a separate authorization lesson

### Lesson 6.3 Completion Notes

#### Business Context

MES Lite business endpoints now require callers to authenticate with a valid JWT before creating or reading production-related data.

#### Problem

The auth package could verify JWTs, but the server did not enforce authentication on requests.

#### Design Discussion

The solution adds a small `auth.Middleware` type. It extracts a bearer token from the `Authorization` header, verifies it with `TokenManager`, stores the resulting `Principal` in request context and calls the next handler. Missing or invalid tokens return `401 Unauthorized` before business handlers execute.

This keeps authentication separate from authorization. Lesson 6.3 answers "who are you?" Lesson 6.4 will answer "are you allowed to do this?"

#### Go Concepts

- `http.Handler` middleware chaining
- request-scoped context values
- unexported context key types
- header parsing with `strings`
- early HTTP rejection before handler execution

#### Architecture Concepts

- route-scoped middleware through Fuego `OptionMiddleware`
- public infrastructure endpoints remain unauthenticated
- protected business endpoints require JWT authentication
- authenticated principal is available without global state

#### Implementation

- Added `auth.Middleware`.
- Added `Authenticate` middleware.
- Added `ContextWithPrincipal` and `PrincipalFromContext`.
- Protected employee, product and production-entry routes.
- Kept `/auth/login`, `/health`, `/ready` and `/version` public.

#### Tests

- Added middleware success test proving the principal reaches request context.
- Added middleware tests for missing and invalid tokens.
- Updated production-entry server route test to send a bearer token.
- Added server test proving production registration requires authentication.
- Added server test proving login remains public.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Business handlers were not changed. Authentication belongs to middleware, not to every individual handler.

#### Code Review

An experienced Go engineer would approve the standard-library middleware shape and route-scoped protection. The main remaining gap is authorization: all authenticated users currently have the same access.

#### Exercises

- Add a test proving `/employees` returns `401` without a token.
- Explain why context values are acceptable for request-scoped identity but not business data.
- Trace a request from `Authorization` header to `PrincipalFromContext`.

#### Interview Questions

- What is HTTP middleware in Go?
- Why use an unexported type as a context key?
- What data should and should not be stored in `context.Context`?
- Why separate authentication middleware from authorization checks?

#### Roadmap Update

- Lesson 6.3 completed.
- Current lesson moved to Lesson 6.4.
- HTTP middleware marked complete in the Knowledge Matrix.
- Known technical debt updated: authentication enforcement exists, role-based authorization remains pending.
