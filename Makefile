.SILENT:
.DEFAULT_GOAL: run
.PHONY: build run stop test test-coverage lint

COVER_FILE=cover.out

MOCK_SRC=./internal/transport/rest/handler.go
MOCK_DST=./internal/mocks/mocks.go

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans

rebuild: build
	docker-compose up --remove-orphans --build

stop:
	docker-compose down

lint:
	golangci-lint run

test:
	go test -coverprofile=$(COVER_FILE) -v ./...
	make --silent test-coverage

test-coverage:
	go tool cover -func=cover.out | grep "total"

clean:
	rm -rf ./.bin $(COVER_FILE)

swag:
	swag init -g ./internal/app/app.go

mock-gen:
	mockgen -source=$(MOCK_SRC) -destination=$(MOCK_DST)