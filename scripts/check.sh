#!/bin/bash

set -e

BIN_DIR="./bin"
LINTER_BINARY="${BIN_DIR}/golangci-lint"

echo "Running code checks..."

${LINTER_BINARY} run ./...
go test -v -cover ./...

echo "Everything is clean and all tests passed!"


