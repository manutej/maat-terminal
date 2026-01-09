# MAAT Makefile
# Follows Commandment #9: Terminal Citizenship

.PHONY: all build run test clean fmt lint install deps

# Binary name
BINARY_NAME=maat
BINARY_PATH=./cmd/maat

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build the project
all: deps fmt test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Lint code
lint:
	@echo "Linting code..."
	$(GOVET) ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Install to system (optional)
install: build
	@echo "Installing to /usr/local/bin..."
	cp $(BINARY_NAME) /usr/local/bin/

# Development mode (watch and rebuild)
dev:
	@echo "Starting development mode..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

# Help
help:
	@echo "MAAT - Makefile commands:"
	@echo "  make build          - Build the binary"
	@echo "  make run            - Build and run"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Lint code"
	@echo "  make deps           - Install dependencies"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install to /usr/local/bin"
	@echo "  make dev            - Run in development mode"
