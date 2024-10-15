# Define variables
BINARY_NAME=music-saas
BUILD_DIR=bin
SRC_DIR=cmd/music-saas

# Default target
all: build

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Clean up binaries and build files
clean:
	rm -rf $(BUILD_DIR)

# Help message
help:
	@echo "Makefile for music-saas"
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean up binaries"
	@echo "  make help       - Show this help message"
