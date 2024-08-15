FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ads-service cmd/main.go

# Creating an image from the builder stage with only the binary and files needed for running
FROM alpine:3.20

WORKDIR /app

# Copy binary 
COPY --from=builder /app/ads-service .
# Copy configs
COPY --from=builder /app/config ./config
# Copy .env, which contains configuration type
COPY --from=builder /app/.env .
# Copy migration files
COPY --from=builder /app/migrations ./migrations

CMD ["./ads-service"]