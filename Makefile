.PHONY: build run-stdio run-http test test-http inspector clean install help

# Binary name
BINARY_NAME=tmdb-mcp-server

# Server configuration
HTTP_PORT ?= 8080

# Build variables
GO=go
GOFLAGS=-v

help:
	@echo 'Usage: make [target] [HTTP_PORT=port]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Environment variables:'
	@echo '  HTTP_PORT        Port for HTTP mode (default: 8080)'
	@echo '  TMDB_API_KEY     Your TMDB API key (required)'
	@echo ''
	@echo 'Examples:'
	@echo '  make run-http'
	@echo '  make run-http HTTP_PORT=3000'
	@echo '  HTTP_PORT=9000 make run-http'

build: ## Build the server binary
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

install: ## Install dependencies
	$(GO) mod download
	$(GO) mod tidy

run-stdio: build ## Run in stdio mode (local)
	@echo "Starting TMDB MCP server in stdio mode..."
	@echo "Press Ctrl+C to stop"
	./$(BINARY_NAME) --mode stdio

run-http: build ## Run in HTTP mode (HTTP_PORT=8080)
	@echo "Starting TMDB MCP server in HTTP mode on port $(HTTP_PORT)..."
	@echo "Access at http://localhost:$(HTTP_PORT)"
	@echo "Test with: make test-http HTTP_PORT=$(HTTP_PORT) (in another terminal)"
	@echo "Or use MCP Inspector: make inspector"
	@echo "Press Ctrl+C to stop"
	./$(BINARY_NAME) --mode http --port $(HTTP_PORT)

test: ## Run tests
	$(GO) test -v ./...

test-http: ## Test HTTP mode with curl (HTTP_PORT=8080)
	@echo "Testing HTTP mode endpoints on port $(HTTP_PORT)..."
	@echo ""
	@echo "=== Listing available tools ==="
	@curl -s -X POST http://localhost:$(HTTP_PORT)/sse \
		-H "Content-Type: application/json" \
		-d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | jq -r '.result.tools[].name' || echo "Error: Is the server running? (make run-http HTTP_PORT=$(HTTP_PORT))"
	@echo ""
	@echo "=== Searching for 'Inception' ==="
	@curl -s -X POST http://localhost:$(HTTP_PORT)/sse \
		-H "Content-Type: application/json" \
		-d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search_movies","arguments":{"query":"Inception"}}}' \
		| jq -r '.result.content[0].text | fromjson | .results[0] | "Found: \(.title) (\(.release_date)) - Rating: \(.rating)"' 2>/dev/null || echo "Error calling search_movies"

inspector: ## Open MCP Inspector for testing (requires Node.js)
	@echo "Opening MCP Inspector..."
	@echo "Make sure server is running: make run-http"
	@echo ""
	@which npx > /dev/null || (echo "Error: npx not found. Install Node.js first." && exit 1)
	@npx @modelcontextprotocol/inspector

test-coverage: ## Run tests with coverage
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

clean: ## Clean build artifacts
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out

fmt: ## Format code
	$(GO) fmt ./...

lint: ## Run linter
	golangci-lint run

docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):latest .

docker-run-stdio: docker-build ## Run Docker container in stdio mode
	docker run -e TMDB_API_KEY=$(TMDB_API_KEY) $(BINARY_NAME):latest

docker-run-http: docker-build ## Run Docker container in HTTP mode (HTTP_PORT=8080)
	docker run -p $(HTTP_PORT):$(HTTP_PORT) -e HTTP_PORT=$(HTTP_PORT) -e TMDB_API_KEY=$(TMDB_API_KEY) $(BINARY_NAME):latest

.DEFAULT_GOAL := help