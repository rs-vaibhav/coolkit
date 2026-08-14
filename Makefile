APP_NAME=coolkit
BUILD_DIR=bin
MAIN_PATH=./cmd/server
DATABASE_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help build run test test-short lint fmt vet clean tidy migrate-up migrate-down migrate-create docker-up docker-down docker-build docker-logs seed coverage all

help:
	@echo "Available targets:"
	@echo "  help           - display all available targets with descriptions"
	@echo "  build          - build the application"
	@echo "  run            - run the application"
	@echo "  test           - run tests"
	@echo "  test-short     - run short tests"
	@echo "  lint           - run linter"
	@echo "  fmt            - format code"
	@echo "  vet            - run go vet"
	@echo "  clean          - remove build dir and coverage files"
	@echo "  tidy           - tidy go modules"
	@echo "  migrate-up     - run migrations up"
	@echo "  migrate-down   - run migrations down"
	@echo "  migrate-create - create a new migration"
	@echo "  docker-up      - start docker containers"
	@echo "  docker-down    - stop docker containers"
	@echo "  docker-build   - build docker containers"
	@echo "  docker-logs    - tail docker logs"
	@echo "  seed           - run database seed"
	@echo "  coverage       - run tests with coverage"
	@echo "  all            - format, vet, lint, test, build"

build:
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)/main.go

test:
	go test ./... -v -cover -race

test-short:
	go test ./... -short

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

tidy:
	go mod tidy

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f

seed:
	go run scripts/seed/main.go

coverage:
	go test ./... -v -coverprofile=coverage.out -race
	go tool cover -html=coverage.out

all: fmt vet lint test build
