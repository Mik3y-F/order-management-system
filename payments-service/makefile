# Variables
SERVICE = payments-service
GOSEC = github.com/securego/gosec/v2/cmd/gosec


# Default target executed when no arguments are given to make.
default: lint test security

lint:
	@echo "Running lint checks for $(SERVICE)..."
	golangci-lint run

test:
	@echo "Running tests for $(SERVICE)..."

	@if ! [ -x "$$(command -v gcloud)" ]; then \
		echo "gcloud is not installed. Please install gcloud."; \
		exit 1; \
	fi

	go test -v -race ./...

security:
	@echo "Running security checks for $(SERVICE)..."

	@if ! [ -x "$$(command -v gosec)" ]; then \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...

run:
	@echo "Running $(SERVICE)..."
	go run cmd/server/main.go

.PHONY: lint test security

