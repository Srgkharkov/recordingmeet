# Используем официальный образ golang для сборки
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY app/go.mod app/go.sum ./
RUN go mod download

# Копируем исходный код приложения
COPY app/ .

# Сборка Go-приложения
RUN go build -o build/main .

# Создаем минимальный образ для запуска
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/build/main /main

# Устанавливаем команду по умолчанию для запуска контейнера
ENTRYPOINT ["/main"]

# # Используем базовый образ Alpine
# FROM alpine:latest

# # Устанавливаем зависимости и Xvfb
# RUN apk update && apk add --no-cache \
#     xvfb-run \
#     chromium \
#     chromium-chromedriver \
#     go \
#     bash

# # Копируем ваше приложение
# COPY your-go-app /app/your-go-app

# # Устанавливаем рабочий каталог
# WORKDIR /app

# # Делаем ваше приложение исполняемым
# RUN chmod +x your-go-app

# # Устанавливаем точку входа
# ENTRYPOINT ["xvfb-run", "-a", "/app/your-go-app"]
