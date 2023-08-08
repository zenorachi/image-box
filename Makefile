.SILENT:
.DEFAULT_GOAL: build
.PHONY: build run stop test test-coverage lint

COVER_FILE=cover.out

build:
	docker-compose up --build

run:
	docker-compose up

stop:
	docker-compose down

test:
	go test -coverprofile=$(COVER_FILE) -v ./...
	make --silent test-coverage

test-coverage:
	go tool cover -func=cover.out | grep "total"

clean:
	rm -rf $(COVER_FILE)

lint:
	golangci-lint run
