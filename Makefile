.PHONY: run test build format vet check docker-build docker-run integration-test

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

fmt:
	go fmt ./...

check:
	@echo "Running linter and tests..."
	@./bin/golangci-lint run ./... || (echo "Linting failed!" && exit 1)
	@go test -v -cover ./... || (echo "Tests failed!" && exit 1)
	@echo "Everything is clean and all tests passed!"

integration-test:
	go test -tags=integration -v ./...

docker-build:
	docker build -t hello-world-go .

docker-run:
	docker run -p 8080:8080 hello-world-go

lint:
	./bin/golangci-lint run ./...