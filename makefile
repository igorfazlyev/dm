.PHONY: run build test docker-up docker-down migrate

run:
	go run cmd/api/main.go

build:
	go build -o bin/dental-api cmd/api/main.go

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

migrate:
	go run cmd/api/main.go

deps:
	go mod download
	go mod tidy

.DEFAULT_GOAL := run
