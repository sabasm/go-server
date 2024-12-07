#!/bin/bash

echo "Running code checks..."

./bin/golangci-lint run ./...
if [ $? -ne 0 ]; then
  echo "Linting failed!"
  exit 1
fi

go test -v -cover ./...
if [ $? -ne 0 ]; then
  echo "Tests failed!"
  exit 1
fi

echo "Everything is clean and all tests passed!"


