.SILENT:
.DEFAULT_GOAL: build
.PHONY: build run stop test test-coverage lint

COVER_FILE=cover.out

MOCK_SRC=./internal/transport/rest/handler.go
MOCK_DST=./internal/mocks/mocks.go

build:
	docker-compose up --build

run:
	docker-compose up

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
	rm -rf $(COVER_FILE)

mock-gen:
	mockgen -source=$(MOCK_SRC) -destination=$(MOCK_DST)

