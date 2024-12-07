run:
	go run ./cmd/main.go

test:
	go test -v -cover ./...

build:
	go build -o bin/app ./cmd/main.go

format:
	go fmt ./...

vet:
	go vet ./...

check:
	@echo "Running linter and tests..."
	@./bin/golangci-lint run ./... || (echo "Linting failed!" && exit 1)
	@go test -v -cover ./... || (echo "Tests failed!" && exit 1)
	@echo "Everything is clean and all tests passed!"


