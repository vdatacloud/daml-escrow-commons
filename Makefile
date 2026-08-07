# daml-escrow-commons

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build all packages
	go build ./...

.PHONY: test
test: ## Run all unit tests
	go test ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: tidy
tidy: ## Tidy go.mod/go.sum
	go mod tidy

.PHONY: verify
verify: build vet test ## Full local gate: build, vet, test (run before pushing)
	@echo "All verification checks passed."
