# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/avito-merch ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/avito-merch /app/avito-merch
COPY migrations ./migrations

EXPOSE 8080

CMD ["./avito-merch"]
