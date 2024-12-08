FROM golang:1.23

WORKDIR /app/hotelier_service

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o hotelier_service ./cmd/hotelier_service

CMD ["./hotelier_service"]