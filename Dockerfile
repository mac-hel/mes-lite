# syntax=docker/dockerfile:1

# Use the official Go image only for compilation. Keeping the toolchain out of
# the runtime image makes the final artifact smaller and reduces attack surface.
FROM golang:1.26.5-alpine AS build

WORKDIR /src

# Copy module files first so Docker can cache downloaded dependencies when only
# application source files change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO is disabled so the binaries do not depend on system C libraries at
# runtime. GOOS=linux makes the build target explicit, and -trimpath removes
# local filesystem paths from the compiled binaries for more reproducible builds.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/mes-lite ./cmd/server \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/migrate ./cmd/migrate

# Alpine keeps the runtime image small while still making it easy to add the few
# OS files this service needs, such as CA certificates and timezone data.
FROM alpine:3.22

# Run the service as a dedicated non-root user. If the process is compromised,
# it should not get root privileges inside the container by default.
RUN addgroup -S meslite \
    && adduser -S -D -H -G meslite meslite \
    && apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy only compiled binaries and runtime assets. Source code, tests and build
# tools stay in the build stage and are not shipped in the production layer.
COPY --from=build /out/mes-lite /app/mes-lite
COPY --from=build /out/migrate /app/migrate
COPY migrations /app/migrations

# These defaults are safe operational defaults. Deployment-specific secrets and
# endpoints, such as DATABASE_URL and JWT_SECRET, must still be supplied by the
# runtime environment.
ENV HOST=0.0.0.0 \
    PORT=9090 \
    MIGRATIONS_DIR=/app/migrations \
    LOG_FORMAT=json \
    LOG_LEVEL=info \
    OTEL_TRACES_EXPORTER=none

EXPOSE 9090

USER meslite:meslite

# The server is the default container process. The migration binary is included
# for release jobs that override the entrypoint or command when running schema
# migrations.
ENTRYPOINT ["/app/mes-lite"]
