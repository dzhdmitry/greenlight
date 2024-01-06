.PHONY: run build up ps bash bash-root down migrate

run:
	go run ./cmd/api

build:
	go build -ldflags='-s' -o=./bin/api ./cmd/api

compose-build:
	docker-compose build

up:
	docker-compose up -d

ps:
	docker-compose ps

bash:
	docker-compose exec -u postgres postgres bash

bash-root:
	docker-compose exec -u root postgres bash

down:
	docker-compose down

migrate:
	migrate -path=./migrations -database=postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable up
