SERVICES := $(shell find ./src -mindepth 1 -maxdepth 1 -type d -exec test -f "{}/go.mod" \; -print)

.PHONY: all test lint build

all: test lint

test:
	@echo "🔍 Running tests for all services with go.mod..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $$service && go test ./... -v || exit 1; \
		cd - > /dev/null; \
	done

lint:
	@echo "🔎 Running golangci-lint for all services..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $$service && golangci-lint run ./... --timeout=2m || exit 1; \
		cd - > /dev/null; \
	done

build:
	@echo "🛠️ Building all services with go.mod..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go build ./... || exit 1; \
		cd - > /dev/null; \
	done
