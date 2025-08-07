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

## run: runs app in Docker
.PHONY: run
run:
	@echo "Starting containers..."
	@docker compose up --build -d
