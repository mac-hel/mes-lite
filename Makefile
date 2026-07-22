.PHONY: build test lint clean run help

APP_NAME = mes-lite
CMD_DIR = ./cmd/server
BUILD_DIR = ./bin

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

test: ## Run tests
	go test ./... -v -count=1

test-race: ## Run tests with race detector # -shuffle=on
	go test ./... -v -count=1 -race

test-cover: ## Run tests with coverage
	go test ./... -v -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

install-tools:
	go install github.com/air-verse/air@latest
	echo "install 'golangci-lint' with your package manager, e.g.: pacman -S golangci-lint"

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format Go code
	go fmt ./...

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
