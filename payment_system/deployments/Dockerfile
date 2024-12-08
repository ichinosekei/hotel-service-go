FROM golang:1.23

WORKDIR /app/payment_system

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o payment_system ./cmd/payment_system

CMD ["./payment_system"]