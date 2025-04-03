.PHONY: run build clean

# GO variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get 
GO_MOD=$(GO_CMD) mod
BINARY_NAME=app
AIR_CMD=air

# Docker variables
DOCKER_CMD=docker
DOCKER_COMPOSE_CMD=docker-compose
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_RUN=$(DOCKER_CMD) run
DOCKER_COMPOSE_UP=$(DOCKER_COMPOSE_CMD) up --build

run: 
	$(AIR_CMD)

build:
	$(GO_BUILD) -o $(BINARY_NAME) .

clean: 
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)

docker-build:
	$(DOCKER_BUILD) -t $(BINARY_NAME)

docker-run:
	$(DOCKER_RUN) -p 8080:8080 $(BINARY_NAME)

docker-compose-up:
	$(DOCKER_COMPOSE_UP)

test: 
	$(GO_TEST) ./...
