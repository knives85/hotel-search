GO     ?= go
BINARY ?= hotel-search
PKG    ?= ./...

.PHONY: help build run test fmt vet tidy clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-8s\033[0m %s\n", $$1, $$2}'

build: ## Compile the binary into ./bin
	$(GO) build -o bin/$(BINARY) ./cmd/hotel-search

run: ## Run the HTTP server
	$(GO) run ./cmd/hotel-search

test: ## Run all tests
	$(GO) test $(PKG)

fmt: ## Format the code
	$(GO) fmt $(PKG)

vet: ## Run static checks
	$(GO) vet $(PKG)

tidy: ## Sync go.mod / go.sum
	$(GO) mod tidy

clean: ## Remove build artifacts
	rm -rf bin
