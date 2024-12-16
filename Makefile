SHELL := /bin/bash
GO := go
GOCOVER := $(GO) tool cover
GOLINT := ./bin/golangci-lint
APP_NAME := hello-world-go
VERSION := v1.0.0
BINARY_DIR := bin
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
MVP_DIR := ./cmd/mvp-api
MIN_COVERAGE := 80

.PHONY: all setup check build lint fmt vet test integration-test docker-build docker-run clean release dev help qcheck check-mvp

all: setup check build

setup:
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

check: fmt lint vet test

check-mvp: setup
	@echo "Running MVP API checks..."
	@echo "Formatting..."
	@$(GO) fmt $(MVP_DIR)/...
	@echo "Linting..."
	@$(GOLINT) run $(MVP_DIR)/...
	@echo "Vetting..."
	@$(GO) vet $(MVP_DIR)/...
	@echo "Testing with coverage..."
	@$(GO) test -v -race -coverprofile=coverage.mvp.out $(MVP_DIR)/...
	@$(GOCOVER) -func=coverage.mvp.out | tee coverage.mvp.txt
	@echo "Verifying coverage threshold..."
	@if [ $$(tail -n 1 coverage.mvp.txt | awk '{print $$NF}' | sed 's/%//') -lt $(MIN_COVERAGE) ]; then \
		echo "Test coverage below $(MIN_COVERAGE)%"; \
		exit 1; \
	fi
	@$(GOCOVER) -html=coverage.mvp.out -o coverage.mvp.html
	@echo "Running integration tests..."
	@$(GO) test -tags=integration -v $(MVP_DIR)/...
	@echo "Building MVP API..."
	@CGO_ENABLED=0 $(GO) build -v -o $(BINARY_DIR)/mvp-api $(MVP_DIR)
	@echo "MVP API checks completed"

lint:
	@$(GOLINT) run ./...

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

test:
	@$(GO) test -v -race -p 1 -coverprofile=coverage.out ./...
	@$(GOCOVER) -func=coverage.out
	@$(GOCOVER) -html=coverage.out -o coverage.html

build:
	@CGO_ENABLED=0 $(GO) build -v -o $(BINARY_DIR)/$(APP_NAME) ./cmd/server

integration-test:
	@$(GO) test -tags=integration -v ./...

docker-build:
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

clean:
	@echo "Cleaning build artifacts..."
	@rm -f main main_test.go cmd/main.go cmd/main_test.go
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html coverage.mvp.out coverage.mvp.html coverage.mvp.txt
	@echo "Clean completed"

release:
	@git checkout -b release/$(VERSION)
	@git add .
	@git commit -m "Release $(VERSION)"
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@git checkout main
	@git merge release/$(VERSION)
	@git push origin main

dev:
	@git checkout -b dev
	@$(MAKE) setup
	@$(MAKE) check

help:
	@echo "Available targets:"
	@echo "  setup            - Set up development environment"
	@echo "  check            - Run linting, vetting, and tests"
	@echo "  check-mvp        - Run MVP API specific checks"
	@echo "  build            - Build the application"
	@echo "  test             - Run tests with coverage"
	@echo "  integration-test - Run integration tests"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  clean            - Remove build artifacts"
	@echo "  release          - Create a new release"
	@echo "  dev              - Set up development environment"
	@echo "  qcheck           - Run tests for a specific directory"

qcheck:
	@if [ -n "$(DIR)" ]; then \
		echo "Running tests for ./cmd/$(DIR)..."; \
		$(GO) test -v ./cmd/$(DIR)...; \
	else \
		echo "Running tests for all packages..."; \
		$(GO) test -v ./...; \
	fi