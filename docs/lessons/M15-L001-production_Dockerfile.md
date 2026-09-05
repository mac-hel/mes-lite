### Lesson 15.1 Scope

Create or review the production Docker image so MES Lite can be built and run as a deployable Go service artifact.

#### Business Context

The application is only production-ready if it can be packaged reproducibly. A reliable Docker image is the deployment unit most teams would use for this kind of API service.

#### Problem

The project has Docker development support, but production readiness requires reviewing whether the runtime image is small, predictable, non-root where practical and able to run the compiled server binary with the required configuration.

#### Design Discussion

L15.1 should focus on the deployable artifact, not the whole release process. Build metadata, CI hardening and smoke tests belong to later lessons so Docker image construction stays clear.

#### Go Concepts

- Go binary builds
- static versus dynamic runtime needs
- linker flags introduction if already needed
- environment-based configuration at runtime

#### Architecture Concepts

- deployment artifact as part of production design
- separating build-time and runtime concerns
- avoiding development-only tooling in production images

#### Tests

- Build the production image or validate the Dockerfile.
- Run the container enough to confirm startup behavior or configuration errors are clear.
- Keep normal tests, build, vet and lint passing.

### Lesson 15.1 Completion Notes

#### Business Context

MES Lite now has a production-oriented Docker image definition for packaging the API as a deployable service artifact.

#### Problem

The project had Docker Compose for local PostgreSQL development, but no production image for the Go service itself. That meant the API could be built locally, but not packaged in the same repeatable shape a deployment pipeline would normally promote.

#### Design Discussion

Added a multi-stage Dockerfile. The build stage uses the Go toolchain and module cache to compile the server and migration binaries. The runtime stage is a smaller Alpine image that contains only runtime assets: binaries, migrations, CA certificates and timezone data.

The final image runs as an unprivileged `meslite` user. This is a default production-hardening step: if the application is compromised, the process should not have root privileges inside the container unless it has a concrete need.

The server remains configured through environment variables. The image sets safe runtime defaults for host, port, migrations directory, log format, log level and trace exporter, while deployment-specific values such as `DATABASE_URL`, `JWT_SECRET` and bootstrap credentials must still be supplied by the runtime environment.

#### Implementation

- Added `Dockerfile` with separate build and runtime stages.
- Built static Linux binaries with `CGO_ENABLED=0` and `-trimpath`.
- Included both `/app/mes-lite` and `/app/migrate` in the image.
- Included `/app/migrations` so the migration command can run from the same artifact.
- Added CA certificates and timezone data to the runtime image.
- Added a non-root `meslite` runtime user.
- Added `.dockerignore` to keep Git metadata, local environment files, coverage output and build artifacts out of the Docker build context.
- Added `make docker-build` as the reproducible local image build command.

#### Tests

- Verified the production image builds with `docker build -t mes-lite:local .`.
- Verified container startup with `docker run --rm mes-lite:local`; without a container-networked PostgreSQL URL, the server exits with a clear structured `ping database` error.
- Verified the runtime container uses a non-root user with UID `100` and GID `101`.
- Verified `/app/mes-lite`, `/app/migrate` and `/app/migrations` exist in the runtime image.
- Verified the Go binaries still compile with `go build ./...`.
- Verified correctness tests with `go test ./... -count=1`.
- Verified static analysis with `go vet ./...`.
- Verified linting with `golangci-lint run ./...`.

#### Refactoring

No application runtime behavior changed. The lesson only added deployment packaging and a Makefile target.

#### Code Review

An experienced Go engineer would approve the Dockerfile direction for this project stage: build tools are excluded from the runtime layer, the container runs as non-root, configuration remains environment-driven and migrations are available in the same artifact.

The image build and basic runtime checks pass locally. A full container smoke test against PostgreSQL belongs to the later CI/CD and smoke-test lesson.

#### Exercises

- Run `make docker-build` on a machine with Docker available.
- Run the container with no `JWT_SECRET` and confirm startup fails with a clear configuration error.
- Compare Alpine runtime with a distroless runtime image and identify the operational trade-offs.

#### Interview Questions

- Why use a multi-stage Docker build for Go services?
- Why is a static Go binary easy to package in a small runtime image?
- Why should production containers usually avoid running as root?
- What belongs in the image versus what belongs in runtime environment configuration?

#### Roadmap Update

- Lesson 15.1 completed.
- Current lesson moved to Lesson 15.2.
- Milestone 15 remains in progress.
