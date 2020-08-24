PROJECT=bot

all: stop build start

build-local:
	go build -o bin/main main.go

build:
	@docker-compose -f docker-compose.yml -p "${PROJECT}" build

stop:
	@docker-compose -p "${PROJECT}" down --rmi local

start:
	@docker-compose -p "${PROJECT}" up

