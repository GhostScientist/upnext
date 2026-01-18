.PHONY: build install clean test run

# Binary name
BINARY=upnext

# Build directory
BUILD_DIR=./bin

# Install directory (respects PREFIX, defaults to /usr/local)
PREFIX ?= /usr/local
INSTALL_DIR=$(PREFIX)/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY) ./cmd/upnext

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
test:
	$(GOTEST) -v ./...

# Install the application
install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)

# Uninstall the application
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run the application
run: build
	$(BUILD_DIR)/$(BINARY)

# Development: run with go run
dev:
	$(GOCMD) run ./cmd/upnext
