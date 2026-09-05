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
version_response=$(curl -fsS "http://127.0.0.1:$HOST_PORT/version")
curl -fsS "http://127.0.0.1:$HOST_PORT/metrics" >/dev/null

if [ "${EXPECTED_VERSION:-}" != "" ]; then
	case "$version_response" in
	*"\"version\":\"$EXPECTED_VERSION\""*) ;;
	*) echo "expected /version to contain version $EXPECTED_VERSION, got: $version_response" >&2; exit 1 ;;
	esac
fi

if [ "${EXPECTED_COMMIT:-}" != "" ]; then
	case "$version_response" in
	*"\"commit\":\"$EXPECTED_COMMIT\""*) ;;
	*) echo "expected /version to contain commit $EXPECTED_COMMIT, got: $version_response" >&2; exit 1 ;;
	esac
fi

if [ "${EXPECTED_BUILD_TIME:-}" != "" ]; then
	case "$version_response" in
	*"\"buildTime\":\"$EXPECTED_BUILD_TIME\""*) ;;
	*) echo "expected /version to contain build time $EXPECTED_BUILD_TIME, got: $version_response" >&2; exit 1 ;;
	esac
fi

docker logs "$CONTAINER_NAME" >/dev/null
