SHELL := /bin/bash
GO := go
GOCOVER := $(GO) tool cover
GOLINT := ./bin/golangci-lint
APP_NAME := hello-world-go
VERSION := v1.0.0
BINARY_DIR := bin
CMD_DIR := ./cmd/server
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

.PHONY: all setup check build lint fmt vet test integration-test docker-build docker-run clean release dev help

all: setup check build

setup:
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

check: lint vet fmt test

lint:
	@$(GOLINT) run ./...

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

test:
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@$(GOCOVER) -func=coverage.out
	@$(GOCOVER) -html=coverage.out -o coverage.html

build:
	@$(GO) build -v -o $(BINARY_DIR)/$(APP_NAME) $(CMD_DIR)

integration-test:
	@$(GO) test -tags=integration -v ./...

docker-build:
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

clean:
	@rm -rf $(BINARY_DIR) coverage.out coverage.html

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
	@echo "  build            - Build the application"
	@echo "  test             - Run tests with coverage"
	@echo "  integration-test - Run integration tests"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  clean            - Remove build artifacts"
	@echo "  release          - Create a new release"
	@echo "  dev              - Set up development environment"

