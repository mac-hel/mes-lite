### Lesson 6.1 Scope

Introduce the authentication vertical slice without protecting business endpoints yet.

#### Business Context

MES Lite now stores real production data. Before role-based access can exist, the system needs a way to represent application users and verify credentials.

#### Problem

The application has employees, products and production entries, but no security identity. Anyone who can reach the API can call every endpoint.

#### Design Discussion

Authentication starts with a small `auth` slice. A user is not the same as an employee: an employee represents a production worker in the business domain, while an auth user represents someone allowed to access the application.

This lesson verifies email/password credentials and returns an opaque access token. JWT signing, request middleware and authorization are intentionally postponed to keep the lesson focused on identity, password hashing and credential errors.

#### Go Concepts

- custom string role type
- password hashing with explicit error handling
- `crypto/rand` for unpredictable token bytes
- authentication errors translated to HTTP `401 Unauthorized`

#### Architecture Concepts

- authentication as its own vertical slice
- separating business employees from security users
- service boundary for credential verification
- transport response hides password hashes with `json:"-"`

### Lesson 6.1 Completion Notes

#### Business Context

MES Lite needs application users before it can protect production and master-data endpoints. Security identity is separate from employees because not every employee is necessarily an API user, and not every API user performs production work.

#### Problem

The API had no login flow. Any client could call every endpoint without proving identity.

#### Design Discussion

The first authentication step introduces a small `internal/auth` vertical slice. It verifies credentials using an auth service and returns an access token, while JWT signing, token verification, request context and endpoint protection are postponed to later lessons.

Passwords are hashed with `bcrypt` from `golang.org/x/crypto`. This is an intentional dependency because the Go standard library does not provide a production password-hashing algorithm. Raw password hashing with SHA-256 would be faster for attackers and inappropriate for real authentication.

The running server can create a bootstrap admin only when both `AUTH_BOOTSTRAP_EMAIL` and `AUTH_BOOTSTRAP_PASSWORD` are configured. There is no default admin password.

#### Go Concepts

- custom string type for roles
- constructor validation for security users
- `json:"-"` to prevent password hashes from being serialized
- `crypto/rand` plus base64 encoding for unpredictable opaque token bytes
- sentinel authentication errors translated to HTTP responses

#### Architecture Concepts

- authentication vertical slice added as `internal/auth`
- auth users kept separate from business employees
- auth handler owns HTTP parsing and error translation
- auth service owns credential verification

#### Implementation

- Added `auth.User`, `auth.Role` and role validation.
- Added bcrypt password hashing and password verification.
- Added `auth.Store` and in-memory store for the first lesson.
- Added `auth.Service.Login`.
- Added `POST /auth/login`.
- Added optional bootstrap admin configuration.
- Wired the auth handler into the server composition root.

#### Tests

- Added user password-hashing tests.
- Added login service tests for valid credentials, wrong password and inactive users.
- Added HTTP handler tests for successful and rejected login.
- Updated server route setup tests.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing business endpoint was protected in this lesson. That avoids mixing password handling, JWT parsing and authorization decisions in one step.

#### Code Review

An experienced Go engineer would approve the separation between auth users and employees, the absence of plaintext password storage and the explicit bootstrap configuration. The main caveat is intentional lesson scope: returned tokens are not yet verified by middleware and users are not durable yet.

#### Exercises

- Explain why an employee and an auth user are different domain concepts.
- Add a table test for every supported role.
- Try replacing bcrypt with SHA-256 and explain why that would be weaker for password storage.

#### Interview Questions

- What is the difference between authentication and authorization?
- Why should passwords be hashed with bcrypt/argon2 instead of SHA-256?
- Why should APIs return `401 Unauthorized` for invalid login credentials?
- Why is a default admin password dangerous?
