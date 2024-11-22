# Stage 1 - Build
FROM golang:1.23 AS builder


WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Проверка существования директории configs
RUN ls -l /app/configs

# Проверка структуры исходного кода
# RUN ls -l /app/cmd/hotelier-service

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hotelier-service ./hotel-service/hotelier-service

# Stage 2 - Run
FROM alpine:latest

WORKDIR /root
# Копируем исполняемый файл из стадии сборки
COPY --from=builder /app/hotelier-service .
# Копируем конфигурацию
COPY configs/hotelier-config.yaml ./config.yaml
# Устанавливаем права на выполнение для бинарного файла
RUN chmod +x ./hotelier-service

# Run the binary
CMD ["./hotelier-service", "-config", "./config.yaml"]
