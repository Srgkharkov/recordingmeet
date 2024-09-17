# Используем официальный образ golang для сборки
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код приложения
COPY . .

# Сборка Go-приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o build/main .

# Создаем минимальный образ для запуска
FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# Устанавливаем зависимости и Xvfb
RUN apk update && apk add --no-cache \
    xvfb-run \
    chromium \
    chromium-chromedriver \
    go \
    bash

# Создаем директорию для хранения записей
RUN mkdir -p /records

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/build/main /main

# Устанавливаем команду по умолчанию для запуска контейнера
ENTRYPOINT ["/main"]
