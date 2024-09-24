# Используем официальный образ golang для сборки
FROM golang:1.23 AS builder

RUN apt-get update && apt-get -y upgrade && apt-get -y install gcc g++ ca-certificates chromium xvfb

RUN mkdir -p /records

# RUN mkdir -p /log/screenshots

# Устанавливаем рабочую директорию
WORKDIR /app

# # Копируем go.mod и go.sum для установки зависимостей
# COPY go.mod go.sum ./
# RUN go mod download

# Копируем исходный код приложения
COPY . .
RUN go mod download

# Собираем приложение
RUN go build -o build/rec ./cmd/server

WORKDIR /


# Сборка Go-приложения
# RUN go build -o build/main .
# CGO_ENABLED=0 GOOS=linux

ENTRYPOINT [ "/app/build/rec" ]