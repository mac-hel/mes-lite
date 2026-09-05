# Skip target name (1st arg) from command line arguments
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# Prevent "No rule to make target" error
%:
	@:

.PHONY: build test lint clean run migrate sqlc help docker-build

APP_NAME = mes-lite
CMD_DIR = ./cmd/server
BUILD_DIR = ./out

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

run: ## Run the server (with air if available, otherwise go run)
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		go run $(CMD_DIR); \
	fi

migrate: ## Run database migrations
	go run ./cmd/migrate

sqlc: ## Generate sqlc database code
	sqlc generate

test: ## Run tests (e.g.: internal/csvimport; ...)
	go test ./$(ARGS) -v -count=1

test-race: ## Run tests with race detector # -shuffle=on
	go test ./$(ARGS) -v -count=1 -race

test-cover: ## Run tests with coverage
	go test ./... -v -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

bench:
	go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries$' -benchmem \
	&& GOMAXPROCS=1 go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -benchtime=2s \
	&& GOMAXPROCS=8 go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -benchtime=2s

prof:
	go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -cpuprofile "/tmp/opencode/csvimport_cpu.out" -memprofile "/tmp/opencode/csvimport_mem.out" \
	&& go tool pprof -top "/tmp/opencode/csvimport_cpu.out" \
	&& go tool pprof -top -alloc_space "/tmp/opencode/csvimport_mem.out" \
	&& go tool pprof -top -alloc_objects "/tmp/opencode/csvimport_mem.out"

prof-escape:
	go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -memprofile "/tmp/opencode/csvimport_l143_mem_before.out" \
	&& go test ./internal/csvimport -run '^$' -gcflags='all=-m'
	#'fix formatting only

prof-escape2:
	go tool pprof -list 'ValidateProductionEntries' "/tmp/opencode/csvimport_l143_mem_before.out" \
	&& go test ./internal/csvimport -run '^$' -gcflags='github.com/mac-hel/mes-lite/internal/csvimport=-m=2' \
	&& go test ./internal/csvimport -run '^$' -gcflags='github.com/mac-hel/mes-lite/internal/csvimport=-m=2' 2>&1 | rg 'internal/csvimport/.+((escapes to heap)|(moved to heap)|(leaks to))' \
	&& go build -gcflags='github.com/mac-hel/mes-lite/internal/csvimport=-m=2' ./internal/csvimport 2>&1 | rg 'internal/csvimport/.+((escapes to heap)|(moved to heap)|(leaks to))'

prof-lines:
	go tool pprof -list ValidateProductionEntries

gc-trace:
	GODEBUG=gctrace=1 go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -benchtime=2s

install-tools:
	go install github.com/air-verse/air@latest
	echo "install 'golangci-lint' with your package manager, e.g.: pacman -S golangci-lint"

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format Go code (e.g.: internal/csvimport; ...)
	go fmt ./$(ARGS)

tidy: ## Tidy Go modules
	go mod tidy
	go mod verify

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)/
	rm -f coverage.out coverage.html

docker-up: ## Start development environment
	@if ! docker info > /dev/null 2>&1; then \
		echo "Error: Docker is not running. Start Docker and try again."; \
		exit 1; \
	fi
	docker compose up -d

docker-down: ## Stop development environment
	@if ! docker info > /dev/null 2>&1; then \
		echo "Error: Docker is not running."; \
		exit 1; \
	fi
	docker compose down

docker-build: ## Build the production Docker image
	docker build -t mes-lite:local .

docker-run: ## Run the production Docker image
	docker run --rm mes-lite:local
	#docker run --rm --entrypoint /bin/sh mes-lite:local -c 'id -u && id -g && test -x /app/mes-lite && test -x /app/migrate && test -d /app/migrations'

self-request:
	./bin/self-request.sh

db-sql:
	docker compose exec -T postgres psql -U meslite -d meslite -c "\d production_entries"
	docker compose exec -T postgres psql -U meslite -d meslite -c "select version_id, is_applied from goose_db_version order by id;"

teach-continue:
	opencode -s ses_f8fd5e0d8ffe49SdpPkME9hlNR

# peak hours: 3-6 i 9-12
opencode-ds-r:
	DEEPSEEK_API_KEY=$(DEEPSEEK_API_KEY) opencode --model deepseek/deepseek-v4-pro
opencode-gpt-r:
	opencode --model openai/gpt-5.4
opencode-gpt-h:
	opencode --model openai/gpt-5.5
