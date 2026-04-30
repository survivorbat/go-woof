MAKEFLAGS := --no-print-directory --silent

default: help

help:
	@echo "Please use 'make <target>' where <target> is one of"
	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format go code, tidy the go.mod file and run the linter
	@go mod tidy
	@go tool golangci-lint fmt
	@go tool ghokin fmt replace ./integration-tests
	@go tool golangci-lint run --fix

lint: ## Lint the code locally
	@golangci-lint run

test: ## Run unit-tests
	go test ./...

