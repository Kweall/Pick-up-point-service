# ---------------------------
# Запуск базы данных в Docker
# ---------------------------

compose-up:
	docker-compose up -d route256_db

compose-down:
	docker-compose down

compose-stop:
	docker-compose stop route256_db

compose-start:
	docker-compose start route256_db

compose-ps:
	docker-compose ps route256_db

# ---------------------------
# Запуск миграций через Goose
# ---------------------------

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

goose-status:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" status

goose-down:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down

APP_NAME := pvz 

build: lint
	go build -o $(APP_NAME) ./

deps-install:
	go mod tidy

run: build
	./$(APP_NAME)

lint:
	golangci-lint run --issues-exit-code=0

test:
	go test -v -cover ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html


all: build

clean:
	rm -f $(APP_NAME)

.PHONY: build deps-install run lint all clean test coverage
