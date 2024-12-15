# Этап сборки бинарного файла
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Компилируем приложение для x86_64 (по умолчанию)
RUN go build -o notification-service ./cmd/notification_service/main.go

# Финальный образ
FROM debian:bookworm-slim

# Устанавливаем рабочую директорию
WORKDIR /

# Копируем бинарный файл
COPY --from=builder /app/notification-service /notification-service

# Указываем порт
EXPOSE 50051

# Запуск приложения
CMD ["/notification-service"]
