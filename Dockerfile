# Этап 1: Сборка приложения
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарника
RUN go build -o server ./cmd/main.go

# Установка утилиты migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Этап 2: Финальный контейнер
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /app/server .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Копируем папку с миграциями
COPY --from=builder /app/migrations ./migrations

COPY static ./static

EXPOSE 8080

CMD ["./server"]
