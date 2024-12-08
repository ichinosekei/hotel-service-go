FROM golang:1.23

WORKDIR /app/booking_service

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o booking_service ./cmd/booking_service

CMD ["./booking_service"]