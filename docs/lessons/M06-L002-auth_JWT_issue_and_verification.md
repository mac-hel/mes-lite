### Lesson 6.2 Scope

Replace temporary opaque login tokens with signed JWT access tokens and add token verification inside the auth package.

#### Business Context

After login, clients need a portable credential they can send with later API requests. The server must be able to verify that credential without storing every active session in memory.

#### Problem

Lesson 6.1 returned an unpredictable token, but the application had no way to validate token claims. That was enough to teach credential verification, but not enough for middleware or protected endpoints.

#### Design Discussion

JWT is introduced through a concrete `auth.TokenManager`. The token manager signs access tokens with HMAC SHA-256 and verifies token signature, expiry and required claims. A concrete type is simpler than introducing an interface before a second implementation exists.

The JWT secret is required at server startup through `JWT_SECRET`. There is intentionally no default secret because a default signing key would make local convenience look like production security.

#### Go Concepts

- concrete dependency injection
- time-based token expiry
- wrapping security errors without leaking parser internals
- validating required claims before trusting token data

#### Architecture Concepts

- JWT issuing and verification belong to the auth slice
- token claims become the future request principal
- configuration controls secrets, not source code defaults
- middleware is postponed until token verification is independently tested

### Lesson 6.2 Completion Notes

#### Business Context

Users can now receive a signed access token after login. This prepares the API for authenticated requests without requiring server-side session storage.

#### Problem

The previous login token was opaque and not verifiable by the application. Middleware could not safely identify a caller from it.

#### Design Discussion

The auth package now owns a `TokenManager` that issues and verifies JWTs. It uses `github.com/golang-jwt/jwt/v5`, which was already present indirectly and is now a direct dependency because the project uses it directly.

JWTs are signed using `HS256`. This is simple and appropriate for a single service as long as `JWT_SECRET` is strong and private. Asymmetric signing can be revisited later if multiple services need to verify tokens without sharing the signing secret.

#### Go Concepts

- concrete types before interfaces
- `time.Time` and expiry claims
- error wrapping with a stable `ErrInvalidToken`
- table-like security tests through focused token cases

#### Architecture Concepts

- token issuing stays behind the auth service
- token verification returns a small `Principal`
- server composition root owns security configuration
- no endpoint protection until middleware has a clear boundary

#### Implementation

- Added `auth.TokenManager`.
- Added JWT issue and verify behavior.
- Added `auth.Principal` extracted from token claims.
- Changed login to return a signed JWT access token.
- Added `JWT_SECRET` configuration.
- Required `JWT_SECRET` at server startup.
- Promoted `github.com/golang-jwt/jwt/v5` to a direct dependency.

#### Tests

- Added JWT issue-and-verify tests.
- Added wrong-secret rejection test.
- Added short-secret rejection test.
- Updated login service tests to verify returned JWTs.
- Updated server and handler tests to use explicit test token managers.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Removed temporary opaque token generation from `auth.Service`. Token generation is now centralized in `TokenManager`, making the next middleware lesson smaller.

#### Code Review

An experienced Go engineer would approve requiring an explicit JWT secret and keeping token verification inside the auth package. The main remaining gap is expected: tokens are verifiable but not yet enforced on requests.

#### Exercises

- Decode a returned JWT at jwt.io and identify which fields are claims, not secrets.
- Change the signing secret in a test and explain why verification fails.
- Explain why JWT payloads must not contain password hashes or sensitive data.

#### Interview Questions

- What problem does a JWT solve?
- What is the difference between signing and encrypting a JWT?
- Why should the server verify the signing algorithm?
- When would asymmetric JWT signing be preferable to HS256?

#### Roadmap Update

- Lesson 6.2 completed.
- Current lesson moved to Lesson 6.3.
- Known technical debt updated: JWT verification exists, middleware enforcement remains pending.
