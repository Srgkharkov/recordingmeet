# Makefile для управления Docker Compose

# Переменные
COMPOSE_FILE := docker-compose.yml

# Правила
.PHONY: up down restart

# Запуск контейнеров в фоне
up:
	docker-compose -f $(COMPOSE_FILE) up -d

# Остановка и удаление контейнеров
down:
	docker-compose -f $(COMPOSE_FILE) down

# Перезапуск контейнеров
restart: down up

# Запуск приложения
runmaingo:
	cd app &&	go run main.go

# Первичное получение сертификатов
 onlycertbot:
	docker-compose run --rm certbot certonly --webroot -w /var/www/certbot -d rec.srgkharkov.ru

