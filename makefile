# .PHONY: run build clean

# GO variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_MOD=$(GO_CMD) mod
GO_CLEAN=$(GO_CMD) clean
GO_GET=$(GO_CMD) get 
GO_INSTALL=$(GO_CMD) install
BINARY_NAME=app
WASM_NAME=main.wasm
WASM_DIR=static/wasm
TEMPLATE_DIR=templates
STATIC_DIR=static

AIR_CMD=air
GOTESTSUM=gotestsum
TEMPL=templ

# Docker variables
DOCKER_CMD=docker
DOCKER_COMPOSE_CMD=docker-compose
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_RUN=$(DOCKER_CMD) run
DOCKER_COMPOSE_UP=$(DOCKER_COMPOSE_CMD) up --build
DOCKER_COMPOSE_DOWN=$(DOCKER_COMPOSE_CMD) down

.PHONY: all dev build build-wasm test clean install-tools docker-up docker-down

all: dev

dev: install-tools
	$(AIR_CMD)

build: generate
	$(GO_BUILD) -o $(BINARY_NAME) .

build-wasm:
	@echo "Building WebAssembly..."
	@mkdir -p $(WASM_DIR)
	@cd wasm && GOOS=js GOARCH=wasm $(GO_BUILD) -o ../$(WASM_DIR)/$(WASM_NAME) .
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)/
	@echo "WASM build complete"

build-all: generate build build-wasm

generate:
	@echo "Generating Templ files..."
	@$(TEMPL) generate

test: 
	$(GOTESTSUM) --format testname

test-watch:
	$(GOTESTSUM) --format testname --watch

# @$(GO_INSTALL) github.com/tinygo-org/tinygo@latest
install-tools:
	@echo "Installing development tools..."
	@$(GO_INSTALL) github.com/air-verse/air@latest
	@$(GO_INSTALL) gotest.tools/gotestsum@latest
	@$(GO_INSTALL) github.com/a-h/templ/cmd/templ@latest
run: 
	$(AIR_CMD)

clean: 
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)

docker-build:
	$(DOCKER_BUILD) -t $(BINARY_NAME)

docker-run:
	$(DOCKER_RUN) -p 8080:8080 $(BINARY_NAME)

docker-compose:
	$(DOCKER_COMPOSE_UP)
docker-compose:
	$(DOCKER_COMPOSE_DOWN)
test: 
	$(GO_TEST) ./...

clean: 
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf ./bin
	@rm -rf ./tmp
	@rm -rf $(WASM_DIR) 

verify: 
	@echo "Verifying environment..."
	@which go || (echo "ERROR: Go not installed"; exit 1)
	@go version
	@which templ || (echo "ERROR: templ not installed - run 'make install-tools'"; exit 1)
	@echo "Environment OK"

debug-wasm: 
	@echo "=== WASM Debug ==="
	@ls -la $(WASM_DIR)
	@file $(WASM_DIR)/$(WASM_NAME)
	@echo "WASM size: $$(du -h $(WASM_DIR)/$(WASM_NAME) | cut -f1)"