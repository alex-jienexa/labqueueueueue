# Билд стадии
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Финальная стадия
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
# Копируем папку migrations в контейнер
COPY --from=builder /app/migrations /app/migrations
CMD ["./main"]