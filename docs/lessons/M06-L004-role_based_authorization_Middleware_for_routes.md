### Lesson 6.4 Scope

Add role-based authorization checks to protected business routes.

#### Business Context

Different MES Lite users have different responsibilities. A production worker should register completed work, but should not administer employees. A manager should maintain products, while an administrator controls user-sensitive master data.

#### Problem

Lesson 6.3 proved user identity, but every authenticated role could call every protected endpoint. Authentication answered "who are you?" but authorization did not yet answer "are you allowed to do this?"

#### Design Discussion

RBAC (Role-Based Access Control) starts with a simple route permission matrix:

- Admin: all business endpoints.
- Manager: product maintenance, master-data reads and production registration.
- Leader: master-data reads and production registration.
- Worker: production registration only.

This is implemented as route-scoped middleware after authentication. Missing or invalid identity returns `401 Unauthorized`; valid identity with insufficient permissions returns `403 Forbidden`.

#### Go Concepts

- variadic function parameters for allowed roles
- map set pattern with `map[Role]struct{}`
- middleware ordering
- HTTP status semantics: `401` vs `403`

#### Architecture Concepts

- authorization policy belongs near route composition for now
- authentication and authorization remain separate middleware steps
- RBAC is explicit instead of hidden in handlers
- route-level policies are easy to review during code review

### Lesson 6.4 Completion Notes

#### Business Context

MES Lite now distinguishes what authenticated users are allowed to do based on their role.

#### Problem

All authenticated users previously had the same permissions. That would let workers administer employees or products, which does not match business responsibilities.

#### Design Discussion

Authorization is implemented as `RequireRole`, a middleware factory that receives allowed roles and checks the authenticated principal already placed in context by `Authenticate`.

The policy is intentionally route-level and explicit. A more complex permission system, policy engine or database-backed permissions table would be premature. The current business rules are simple enough for direct route composition.

#### Go Concepts

- variadic parameters with `roles ...Role`
- efficient membership checks with a map set
- middleware ordering with authentication before authorization
- stable HTTP semantics for security failures

#### Architecture Concepts

- RBAC at the HTTP boundary
- explicit route permission matrix
- separation between identity proof and permission checks
- handlers remain free of security boilerplate

#### Implementation

- Added `auth.Middleware.RequireRole`.
- Added `403 Forbidden` response for insufficient permissions.
- Applied admin-only permissions to employee mutation routes.
- Allowed admins, managers and leaders to read master data.
- Allowed admins and managers to maintain products.
- Allowed admins, managers, leaders and workers to register production.

#### Tests

- Added middleware tests for matching role, missing principal and forbidden role.
- Added server test proving worker cannot create employees.
- Added server test proving leader can list products.
- Added server test proving worker can register production.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No handler changes were required. This confirms the middleware boundary is doing the right job.

#### Code Review

An experienced Go engineer would approve the explicit route policy for this project size. The main follow-up is milestone-level review: decide whether in-memory auth users are acceptable for now or whether durable auth-user persistence must be completed before Milestone 6 closes.

#### Exercises

- Add a test proving a manager can create a product.
- Add a test proving a worker cannot list employees.
- Explain why `403 Forbidden` is different from `401 Unauthorized`.

#### Interview Questions

- What is RBAC?
- Why should authentication and authorization be separate concepts?
- When would route-level authorization become insufficient?
- How would you model permissions if roles became too coarse?

#### Roadmap Update

- Lesson 6.4 completed.
- Current lesson moved to Lesson 6.5.
- Known technical debt updated: role checks exist; durable auth-user persistence remains under review.
