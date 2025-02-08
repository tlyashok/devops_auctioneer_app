DB_CONTAINER_NAME=auction_db
APP_CONTAINER_NAME=auction_app

-include .env

# Создание .env файла
env:
	@echo Creating .env file with environment variables...
	@echo DB_HOST=db > .env
	@echo DB_PORT=5433 >> .env
	@echo DB_USER=auction_user >> .env
	@echo DB_PASSWORD=auction_password >> .env
	@echo DB_NAME=auction_db >> .env
	@echo JWT_SECRET_KEY=secret >> .env
	@echo APP_PORT=8000 >> .env


# Сборка и запуск контейнеров
build:
	@echo "Building Docker images..."
	docker-compose build

up:
	@echo "Starting services..."
	docker-compose up -d --build --remove-orphans

down:
	@echo "Stopping and removing containers..."
	docker-compose down --remove-orphans

restart:
	@echo "Restarting application..."
	$(MAKE) down
	$(MAKE) up

logs:
	@echo "Showing application logs..."
	docker logs -f $(APP_CONTAINER_NAME)

# Работа с БД
open_db:
	@echo "Подключаемся к базе данных в контейнере..."
	@docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)

migrate_up:
	@echo "Applying database migrations..."
	docker exec -it $(APP_CONTAINER_NAME) migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@db:${DB_PORT}/${DB_NAME}?sslmode=disable" up

migrate_down:
	@echo "Reverting last database migration..."
	docker exec -it $(APP_CONTAINER_NAME) migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@db:${DB_PORT}/${DB_NAME}?sslmode=disable" down

migrate_force:
	@echo "Forcing migration version..."
	docker exec -it $(APP_CONTAINER_NAME) migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@db:${DB_PORT}/${DB_NAME}?sslmode=disable" force

# Локальный запуск (без Docker)
run:
	@echo "Building and running Go application..."
	go run cmd/main.go

build_go:
	@echo "Building Go binary..."
	go build -o server ./cmd/main.go
