# Project settings
APP_BACKEND := backend
APP_MIDDLEWARE := middleware
APP_SERVICE := service
APP_CLIENT := client

BIN_DIR := bin

# Go tools
GO := go
BUF := buf

# Output binaries
BACKEND_BIN := $(BIN_DIR)/backend
MIDDLEWARE_BIN := $(BIN_DIR)/middleware
SERVICE_BIN := $(BIN_DIR)/service
CLIENT_BIN := $(BIN_DIR)/client

# Default target
all: build

## --- Build targets ---
build: proto $(BACKEND_BIN) $(MIDDLEWARE_BIN) $(SERVICE_BIN) $(CLEINT_BIN)

$(BACKEND_BIN): $(APP_BACKEND)/main.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ ./$(APP_BACKEND)

$(MIDDLEWARE_BIN): $(APP_MIDDLEWARE)/main.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ ./$(APP_MIDDLEWARE)

$(SERVICE_BIN): $(APP_SERVICE)/main.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ ./$(APP_SERVICE)

$(CLEINT_BIN): $(APP_CLIENT)/main.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ ./$(APP_CLIENT)

## --- Run executables ---
run-client: $(CLEINT_BIN)
	$(CLEINT_BIN)
	
run-backend: $(BACKEND_BIN)
	$(BACKEND_BIN)

run-middleware: $(MIDDLEWARE_BIN)
	$(MIDDLEWARE_BIN)

run-service: $(SERVICE_BIN)
	$(SERVICE_BIN)

## --- Protobuf generation ---
proto:
	$(BUF) generate

lint-proto:
	$(BUF) lint

## --- Testing ---
test:
	$(GO) test ./... -v

## --- Clean ---
clean:
	rm -rf $(BIN_DIR)

.PHONY: all build run-backend run-middleware run-service run-client proto lint-proto test clean
