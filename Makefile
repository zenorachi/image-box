.SILENT:
.DEFAULT_GOAL: build
.PHONY: build run stop lint

build: stop
	docker-compose up --build

run: stop
	docker-compose up

stop:
	docker-compose down

lint:
	golangci-lint run