FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o auth-service cmd/main.go

# Копируем с прошлого билда только используемые файлы и бинарник
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /app/.env .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

CMD ["./auth-service"]