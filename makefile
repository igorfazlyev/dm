.PHONY: help dev build up down logs clean

help:
	@echo "Available commands:"
	@echo "  make dev      - Start development environment"
	@echo "  make build    - Build all containers"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - Follow logs"
	@echo "  make clean    - Clean up containers and volumes"

dev:
	docker-compose up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	rm -rf backend/tmp
	rm -rf frontend/node_modules
	rm -rf frontend/dist

# Backend specific
backend-shell:
	docker-compose exec backend sh

backend-migrate:
	docker-compose exec backend go run cmd/api/main.go migrate

# Frontend specific
frontend-shell:
	docker-compose exec frontend sh

frontend-install:
	cd frontend && npm install
