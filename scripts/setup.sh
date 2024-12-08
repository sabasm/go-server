#!/bin/bash

set -e

BIN_DIR="./bin"
LINTER_BINARY="${BIN_DIR}/golangci-lint"

mkdir -p ${BIN_DIR}

if ! [ -x "$(command -v ${LINTER_BINARY})" ]; then
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
fi

go mod tidy

echo "Setup completed successfully."


