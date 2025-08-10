## help: prints this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## test: runs tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

## load-test: runs load test
.PHONY: load-test
load-test:
	@echo "Running load test..."
	@docker run --rm -i grafana/k6:1.1.0 run - <./docker/k6/load-test.js

## run: runs app in Docker
.PHONY: run
run:
	@echo "Starting containers..."
	@docker compose up --build -d

## swagger: generates Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go