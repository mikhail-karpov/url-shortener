PHONY: test
make test:
	@echo "Running tests..."
	@go test -v ./...

PHONY: run
make run:
	@echo "Running app..."
	@go run ./cmd