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
