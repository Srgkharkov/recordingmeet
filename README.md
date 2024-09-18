# Meeting Recorder

## Описание

**Meeting Recorder** — это веб-сервис для записи видеоконференций из различных сервисов (Google Meet, Zoom) с использованием Headless Chrome через библиотеку `chromedp`. Сервис может сохранять записи встреч в локальные директории, архивировать их, а также предоставлять возможность скачивания архивов с видеозаписями через веб-интерфейс.

## Структура проекта

Проект организован в соответствии с рекомендациями по лучшим практикам для Go:

/app /cmd /recorder main.go # Точка входа в приложение /internal /auth # Пакет для работы с JWT auth.go # Функции для создания и проверки JWT /handlers # Пакет для HTTP-обработчиков download_handler.go # Обработчик для скачивания архивов record_handler.go # Обработчик для начала записи конференции /services # Логика работы с конференциями meet_service.go # Взаимодействие с Google Meet /utils # Вспомогательные функции file_utils.go # Работа с файловой системой chromedp_utils.go # Вспомогательные функции для chromedp /mediarecorderjs # JavaScript для запуска MediaRecorder в браузере go.mod # Модуль Go Dockerfile # Dockerfile для сборки приложения /docker-compose.yml # Компоновка сервисов (сервер, NGINX, Certbot) /log # Директория для хранения логов и скриншотов ошибок /records # Директория для сохранения записей встреч


## Как запустить проект

### 1. Зависимости

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### 2. Переменные окружения

Создайте файл `.env` в корне проекта и добавьте туда переменные окружения:

JWT_SECRET_KEY=your_secret_key


### 3. Запуск

1. Соберите контейнеры и запустите приложение:

```bash
docker-compose up --build

2. Сервер будет доступен по адресу http://localhost:8080.

### 4. Использование

Авторизация
Для работы API необходимо передавать валидный JWT токен в заголовке Authorization:

git clone https://github.com/Srgkharkov/recordingmeet.git

Authorization: Bearer <ваш-токен>

Эндпоинты
Начало записи конференции

Метод: GET
URL: /record?link=<ссылка_на_конференцию>
Пример: http://localhost:8080/record?link=https://meet.google.com/abc-defg-hij
Запускает запись конференции. Запись сохраняется в директории /records.

Скачивание записей

Метод: GET
URL: /download?recordsid=<идентификатор_записи>
Пример: http://localhost:8080/download?recordsid=GM_abc-defg-hij_1234567890
Скачивает архив с файлами записей конференции.
### 5. Тестирование
Для тестирования API можно использовать любой HTTP-клиент, например, Postman или curl.

Пример запроса для начала записи:

curl -X GET "http://localhost:8080/record?link=https://meet.google.com/abc-defg-hij" -H "Authorization: Bearer <ваш-токен>"

Пример запроса для скачивания архива:

curl -X GET "http://localhost:8080/download?recordsid=GM_abc-defg-hij_1234567890" -H "Authorization: Bearer <ваш-токен>" -o output.zip

Docker
В проекте используется docker-compose.yml, который собирает три контейнера:

Приложение recorder
nginx для проксирования запросов и работы с SSL
certbot для автоматической генерации сертификатов SSL
Построение контейнеров

docker-compose build

Запуск контейнеров

docker-compose up


Логи и записи
Логи приложения и снимки экрана при ошибках сохраняются в директорию /log.
Все записи конференций сохраняются в директорию /records.
Примечания
Для работы с Google Meet приложение использует chromedp в режиме headless.
Поддержка Zoom может быть добавлена позже, в данный момент реализована только интеграция с Google Meet.
Лицензия
Проект распространяется под лицензией MIT.


Этот `README.md` описывает ключевые компоненты проекта, способ его запуска и использования, а также некоторые важные моменты, связанные с тестированием и структурой проекта.
