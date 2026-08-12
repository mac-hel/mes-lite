### Lesson 6.5 Scope

Harden authentication for milestone completion by making auth users durable, documenting OpenAPI bearer security and reviewing the milestone.

#### Business Context

Authentication is not useful if users disappear after an application restart. The first production-safe identity path must persist at least the bootstrap administrator used to access the secured API.

#### Problem

Auth users were still in-memory while the rest of the business data was PostgreSQL-backed. This meant login worked only until restart, and Milestone 6 would have ended with a major security usability gap.

#### Design Discussion

The project now persists auth users in PostgreSQL through the auth vertical slice. Full user-management CRUD is intentionally postponed. The current business requirement is secure application access, and the minimum durable path is an idempotent bootstrap admin plus database-backed login.

OpenAPI now advertises a `bearerAuth` security scheme and marks protected routes as requiring JWT bearer authentication.

#### Go Concepts

- persistence adapter reuse with sqlc
- idempotent startup behavior
- error translation for duplicate users
- direct dependency promotion when a previously indirect package is imported

#### Architecture Concepts

- durable security identity
- bootstrap workflow instead of default credentials
- OpenAPI as part of the security contract
- concrete persistence before full management workflows

### Lesson 6.5 Completion Notes

#### Business Context

MES Lite now has a restart-safe authentication path. A configured bootstrap administrator is stored in PostgreSQL and can log in after restarts.

#### Problem

In-memory auth users were acceptable for learning login mechanics but not acceptable for completing an authentication milestone.

#### Design Discussion

Auth-user persistence was added without creating full user-management CRUD. This is a deliberate minimal design: CRUD would introduce extra authorization and lifecycle rules that the business has not needed yet.

Bootstrap admin creation is idempotent. If the configured email already exists, startup leaves the existing user unchanged instead of overwriting the password on every restart.

#### Go Concepts

- sqlc-generated auth queries
- PostgreSQL `bytea` for password hashes
- idempotent bootstrap with `EXISTS`
- domain error mapping from PostgreSQL constraint errors

#### Architecture Concepts

- auth vertical slice owns auth persistence
- generated auth database types stay below the auth package boundary
- server composition root selects durable auth store
- OpenAPI security scheme documents runtime security

#### Implementation

- Added migration `0005_create_auth_users.sql`.
- Added auth sqlc queries and generated `internal/auth/authdb`.
- Added `auth.PostgresStore`.
- Added idempotent `EnsureBootstrapAdmin`.
- Switched `cmd/server` from in-memory auth store to PostgreSQL auth store.
- Added OpenAPI `bearerAuth` security scheme.
- Added bearer security requirements to protected business routes.

#### Tests

- Added auth PostgreSQL store integration tests.
- Tested save/find by email.
- Tested duplicate email maps to `ErrAlreadyExists`.
- Tested missing email maps to `ErrNotFound`.
- Tested bootstrap admin creation is idempotent and does not overwrite an existing password.
- Verified with `sqlc generate`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The server now uses PostgreSQL-backed auth users. In-memory auth store remains for fast tests, matching the existing employee/product testing pattern.

#### Code Review

An experienced Go engineer would approve completing durable auth persistence before closing the milestone. The remaining limitation is deliberate: there is no user-management CRUD yet. That should be introduced only when the administrator workflow is defined.

#### Exercises

- Add a test proving invalid role strings are rejected by the database.
- Explain why bootstrap should not overwrite an existing admin password.
- Inspect the OpenAPI document and identify `bearerAuth`.

#### Interview Questions

- Why is idempotent startup important?
- Why should password hashes be stored as bytes instead of plaintext?
- Why keep generated sqlc auth types out of HTTP handlers?
- What risks appear when adding user-management CRUD?

#### Roadmap Update

- Lesson 6.5 completed.
- Milestone 6 completed.
- Current milestone moved to Milestone 7.
- Known technical debt updated: durable auth-user persistence completed; full auth-user management postponed.
