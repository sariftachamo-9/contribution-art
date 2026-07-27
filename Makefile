# Contribution Art (sarif) Makefile

BINARY_NAME=sarif
GO=go

.PHONY: all build run test fmt clean help

all: build

## build: Build the Go CLI binary
build:
	@echo "Building binary $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) main.go

## run: Run the application with default dry-run flag
run:
	$(GO) run main.go -pattern "SARIF" -dry-run

## test: Run unit and integration tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

## fmt: Format source code using go fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## clean: Remove built binary artifacts
clean:
	@echo "Cleaning build artifacts..."
	@if exist $(BINARY_NAME) del $(BINARY_NAME)
	@if exist $(BINARY_NAME).exe del $(BINARY_NAME).exe

## help: Display available Makefile targets
help:
	@echo Contribution Art Build Tasks:
	@echo.
	@echo   make build   - Compile the binary ($(BINARY_NAME))
	@echo   make run     - Run CLI preview mode (SARIF)
	@echo   make test    - Run unit tests
	@echo   make fmt     - Format source files
	@echo   make clean   - Remove compiled binaries
