BINARY_NAME=server
SRC_DIR=.

default: build

install:
	@echo "Installing service on systemd"
	@cp -f server.service /etc/systemd/system/
	@systemctl enable server.service
	@systemctl start server.service
	@echo "Service installed"

build:
	@echo "Building the binary..."
	@go build -o $(BINARY_NAME) $(SRC_DIR)
	@echo "Build complete. Binary is located at ./$(BINARY_NAME)"

run: build
	@echo "Running the application..."
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning the build artifacts..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete."

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "Dependencies installed."

fmt:
	@echo "Formatting the code..."
	@go fmt ./...
	@echo "Code formatted."

lint:
	@echo "Linting the code..."
	@golangci-lint run
	@echo "Lint complete."

.PHONY: build run clean deps fmt lint
