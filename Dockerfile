FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Перейти в папку с исходным кодом
WORKDIR /app/cmd/anti_spam_bot

# Собрать бинарник
RUN go build -o myapp

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY /.env /app

# Копировать бинарник из сборочного этапа
COPY --from=builder /app/cmd/anti_spam_bot/myapp .

CMD ["./myapp"]