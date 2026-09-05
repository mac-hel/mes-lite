#!/bin/sh

set -eu

IMAGE=${IMAGE:-mes-lite:local}
CONTAINER_NAME=${CONTAINER_NAME:-mes-lite-smoke}
HOST_PORT=${HOST_PORT:-9090}
DATABASE_URL=${DATABASE_URL:-postgres://meslite:meslite@host.docker.internal:5432/meslite?sslmode=disable}
JWT_SECRET=${JWT_SECRET:-smoke-test-jwt-secret-with-at-least-32-characters}

cleanup() {
	docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
}

trap cleanup EXIT INT TERM
cleanup

docker run --rm \
	--add-host=host.docker.internal:host-gateway \
	-e DATABASE_URL="$DATABASE_URL" \
	--entrypoint /app/migrate \
	"$IMAGE"

docker run -d \
	--name "$CONTAINER_NAME" \
	--add-host=host.docker.internal:host-gateway \
	-p "127.0.0.1:$HOST_PORT:9090" \
	-e DATABASE_URL="$DATABASE_URL" \
	-e JWT_SECRET="$JWT_SECRET" \
	"$IMAGE" >/dev/null

for _ in $(seq 1 30); do
	if curl -fsS "http://127.0.0.1:$HOST_PORT/ready" >/dev/null 2>&1; then
		break
	fi
	sleep 1
done

curl -fsS "http://127.0.0.1:$HOST_PORT/health" >/dev/null
curl -fsS "http://127.0.0.1:$HOST_PORT/ready" >/dev/null
curl -fsS "http://127.0.0.1:$HOST_PORT/version" >/dev/null
curl -fsS "http://127.0.0.1:$HOST_PORT/metrics" >/dev/null

docker logs "$CONTAINER_NAME" >/dev/null
