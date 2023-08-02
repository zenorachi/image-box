.SILENT:

APP = image-box-app

build:
	docker-compose up --build $(APP)

run:
	docker-compose up $(APP)

stop:
	docker-compose down

lint:
	golangci-lint run