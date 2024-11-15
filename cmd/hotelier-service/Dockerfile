# Stage 1 - Build
FROM golang:1.18 AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o hotelier-service ./cmd/hotelier-service

# Stage 2 - Run
FROM alpine:latest

WORKDIR /root
COPY --from=builder /app/hotelier-service .
COPY configs/hotelier-config.yaml ./config.yaml

# Run the binary
CMD ["./hotelier-service", "-config", "./config.yaml"]
