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
 	docker run -it --rm  -v /root/recordingmeet/nginx/certs:/etc/letsencrypt  -v /root/recordingmeet/nginx:/var/lib/letsencrypt  -p 80:80  certbot/certbot certonly --standalone -d rec.srgkharkov.ru
	# docker-compose run --rm certbot certonly --webroot -w /var/www/certbot -d rec.srgkharkov.ru

