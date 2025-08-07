PHONY: test
make test:
	@echo "Running tests..."
	@go test -v ./...

PHONY: run
make run:
	@echo "Running app..."
	@go run ./cmd

PHONY: run-docker
make run-docker:
	@echo "Starting containers..."
	@docker compose up --build -d
