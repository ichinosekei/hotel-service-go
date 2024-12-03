FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o delivery-service ./cmd/delivery_service/main.go

FROM debian:bookworm-slim
WORKDIR /

# Установка корневых сертификатов
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/delivery-service /delivery-service

CMD ["/delivery-service"]