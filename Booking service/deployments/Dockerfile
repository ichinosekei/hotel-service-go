FROM golang:1.23
WORKDIR /app
RUN apt-get update && apt-get install -y netcat-openbsd && apt-get clean
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o booking_service ./cmd/booking_service
EXPOSE 8080
CMD ["./booking_service"]