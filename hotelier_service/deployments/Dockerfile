FROM golang:1.23

WORKDIR /app/hotelier_service

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o hotelier_service ./cmd/hotelier_service
#COPY .env.dev.dev .env.dev.dev
#COPY ../.env.dev.dev /app/.env.dev.dev
#COPY deployments/.env.dev.dev /app/deployments/.env.dev.dev
CMD ["./hotelier_service"]