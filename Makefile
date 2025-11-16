APP_NAME=service
CMD_DIR=./cmd

.PHONY: build docker-up docker-down

build:
	go build -o $(APP_NAME) $(CMD_DIR)/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v