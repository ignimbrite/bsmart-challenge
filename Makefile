APP_NAME ?= bsmart-challenge
DOCKER_COMPOSE ?= $(shell if command -v docker-compose >/dev/null 2>&1; then echo docker-compose; elif docker compose version >/dev/null 2>&1; then echo "docker compose"; else echo ""; fi)

.PHONY: run tidy test lint compose-up compose-down

run:
	go run ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

lint:
	golangci-lint run ./...

compose-up:
ifeq ($(strip $(DOCKER_COMPOSE)),)
	@echo "Docker Compose no encontrado. Instala docker-compose o habilita el plugin 'docker compose'." && exit 1
else
	$(DOCKER_COMPOSE) up -d
endif

compose-down:
ifeq ($(strip $(DOCKER_COMPOSE)),)
	@echo "Docker Compose no encontrado. Instala docker-compose o habilita el plugin 'docker compose'." && exit 1
else
	$(DOCKER_COMPOSE) down
endif

.PHONY: docker-build docker-up

docker-build:
	docker build -t bsmart-api:latest .

docker-up:
	$(DOCKER_COMPOSE) up --build -d
