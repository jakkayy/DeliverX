FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o notification-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/notification-service .
EXPOSE 8086
CMD ["./notification-service"]
