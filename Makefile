.PHONY: build run test clean deps swagger

# Default binary output
BINARY_NAME=validra-engine

# Set the correct module path
MODULE_PATH=github.com/arifsetyawan/validra

# Directories
SRC_DIR=./src
BUILD_DIR=./build
CMD_DIR=./src/cmd
HTTP_CMD = $(CMD_DIR)/http

# Build the application
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

# Run the application
run:
	@go run $(HTTP_CMD)

# Test the application
test:
	@echo "Testing..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f validra.db
	@go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download

# Update swagger documentation
swagger:
	@echo "Updating Swagger documentation..."
	@swag init -g $(HTTP_CMD)/main.go

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code for potential issues
vet:
	@echo "Vetting code..."
	@go vet ./...

# Run migrations
migrate:
	@echo "Running migrations..."
	@go run $(CMD_DIR)/migrate/main.go

# Generate and install dependencies
install: deps build
	@echo "Installing..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Show help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Vet code for potential issues"
	@echo "  make install  - Install the application"
	@echo "  make help     - Show this help message"

# Default target
default: build