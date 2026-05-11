FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o order-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/order-service .
EXPOSE 8083
CMD ["./order-service"]
