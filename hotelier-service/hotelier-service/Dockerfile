# Stage 1 - Build
FROM golang:1.23 AS builder


WORKDIR /app/hotelier-service

# Copy and download dependencies
COPY hotelier-service/go.mod hotelier-service/go.sum ./

RUN go mod download

# Copy the source code
WORKDIR /app
COPY . .

# Перейти в директорию, где находится main.go
WORKDIR /app/hotelier-service/hotelier-service

## Проверка существования директории configs
#RUN #ls -l /app/configs
#
## Проверка структуры исходного кода
## RUN ls -l /app/cmd/hotelier-service
## Проверяем структуру папок (отладка)
#RUN ls -l /app/hotelier-service
#RUN ls -l /app/hotelier-service/cmd



# Build the application
#WORKDIR /app/hotelier-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/hotelier-service-bin .

# Stage 2 - Run
FROM alpine:latest

WORKDIR /root
# Копируем исполняемый файл из стадии сборки
COPY --from=builder /app/hotelier-service-bin .
# Копируем конфигурацию
COPY hotelier-service/configs/hotelier-config.yaml ./config.yaml
# Устанавливаем права на выполнение для бинарного файла
RUN chmod +x ./hotelier-service-bin

# Run the binary
CMD ["./hotelier-service-bin", "-config", "./config.yaml"]
