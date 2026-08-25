# Skip target name (1st arg) from command line arguments
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# Prevent "No rule to make target" error
%:
	@:

.PHONY: build test lint clean run migrate sqlc help

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

self-request:
	./bin/self-request.sh

db-sql:
	docker compose exec -T postgres psql -U meslite -d meslite -c "\d production_entries"
	docker compose exec -T postgres psql -U meslite -d meslite -c "select version_id, is_applied from goose_db_version order by id;"

teach-continue:
	opencode -s ses_fd78892b4ffegwLCmc7owk6Uvy

# peak hours: 3-6 i 9-12
opencode-ds-r:
	DEEPSEEK_API_KEY=$(DEEPSEEK_API_KEY) opencode --model deepseek/deepseek-v4-pro
opencode-gpt-r:
	opencode --model openai/gpt-5.4
opencode-gpt-h:
	opencode --model openai/gpt-5.5
