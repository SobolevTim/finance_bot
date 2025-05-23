# Используем официальный образ Go для сборки
FROM golang:1.23.4-alpine AS builder

# Установим рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код проекта в контейнер
COPY . .

# Компилируем приложение в бинарный файл
RUN go build -o finance-bot cmd/telegram-bot/main.go

# Используем минимальный образ на основе Alpine для уменьшения размера
FROM alpine:latest

# Установим пакет tzdata для временных зон
RUN apk add --no-cache tzdata

# Создаем рабочую директорию для контейнера
WORKDIR /root/

# Копируем скомпилированное приложение из стадии сборки
COPY --from=builder /app/finance-bot .

# Копируем файл .env (если используется)
# COPY internal/pkg/config/default.yaml ./internal/pkg/config/default.yaml

# Команда для запуска приложения
CMD ["./finance-bot"]