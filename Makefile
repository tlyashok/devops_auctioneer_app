DB_CONTAINER_NAME=db

all: build

run:
	@echo "Building the Go application..."
	go run cmd/main.go

up:
	@echo "Starting PostgreSQL container using Docker Compose..."
	docker-compose up -d --build --remove-orphans

env:
	@echo "Creating .env file with environment variables..."
	echo DB_HOST=db > .env
	echo DB_PORT=5432 >> .env
	echo DB_USER=auction_user >> .env
	echo DB_PASSWORD=auction_password >> .env
	echo DB_NAME=auction_db >> .env
	echo JWT_SECRET_KEY=secret >> .env

down:
	@echo "Cleaning up Docker containers and images..."
	docker-compose up -d --remove-orphans

-include .env

open_db:
	@echo "Подключаемся к базе данных в контейнере..."
	@docker exec -it auction_db psql -U $(DB_USER) -d $(DB_NAME)

migrate_up:
	migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

migrate_down:
	migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" down

start: build db run
