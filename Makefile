.SILENT:
.DEFAULT_GOAL: build
.PHONY: build run stop lint

build:
	docker-compose up --build

run:
	docker-compose up

stop:
	docker-compose down

lint:
	golangci-lint run