# Meeting Recorder

## Описание

**Meeting Recorder** — это веб-сервис для записи видеоконференций из различных сервисов (Google Meet, Zoom) с использованием Headless Chrome через библиотеку `chromedp`. Сервис может сохранять записи встреч в локальные директории, архивировать их, а также предоставлять возможность скачивания архивов с видеозаписями через веб-интерфейс.

## Структура проекта

Проект организован в соответствии с рекомендациями по лучшим практикам для Go:

```bash
/recordingmeet
├── /app                              # Директория приложения
│   ├── /cmd                          # Директория для запускаемых приложений
│   │   └── /server                   # Основное приложение сервера
│   │       └── main.go               # Точка входа для запуска сервера
│   ├── /internal                     # Логика, которая используется только внутри проекта
│   │   ├── /auth                     # Пакет для аутентификации
│   │   │   └── jwt.go                # Логика валидации токена
│   │   ├── /handler                  # Пакет для обработки запросов
│   │   │   └── download_handler.go   # Логика обработки HTTP-запросов скачивания архива
│   │   │   └── log_handler.go        # Логика обработки HTTP-запросов скачивания лога
│   │   │   └── record_handler.go     # Логика обработки HTTP-запросов на запись встреч
│   │   ├── /meet                     # Пакет для работы с конференц платформами
│   │   │   ├── meet_service.go       # Парсинг запроса и определение основных параметров записи
│   │   │   └── google_meet. go       # Основная логика по работе с Google Meet
│   │   └── /utils                    # Пакет парсинга запроса
│   │       ├── file_utils.go         # Функции для работы с файлами
│   │       └── mediarecorderjs.go    # JS скрипт для встраивания на страницу, который находит и записывает аудио и видеопотоки
│   ├── Dockerfile                    # Dockerfile для сборки образа
│   ├── go.mod                        # Файл зависимостей Go
│   └── go.sum                        # Контрольные суммы зависимостей
├── /nginx                            # Директория nginx файлов
│   ├── /certs                        # Директория для файлов сертификатов SSL
│   └── nginx.conf                    # Конфигурационный файл nginx
├── docker-compose.yml                # Dockerfile для сборки образа
├── README.md                         # Документация по проекту
└── Makefile                          # Makefile для автоматизации команд
```

## Как запустить проект

### 1. Зависимости

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- Make

### 2. Переменные окружения

Создайте файл `.env` в корне проекта и добавьте туда переменные окружения:

```bash
JWT_SECRET_KEY=your_secret_key
NOTIFICATION_URL=https://
```


### 3. Запуск приложения

1. Выполните клонирование репозитория:

```bash
git clone https://github.com/Srgkharkov/recordingmeet.git
```

2. Соберите образ:

```bash
cd recordingmeet
docker build ./app -t recordingmeet
docker run -p 8080:8080 --env-file .env recordingmeet
```   
3. Запустите контейнер, сервер будет доступен по адресу http://localhost:8080.
```bash
docker run -p 8080:8080 --env-file .env recordingmeet
```   

### 4. Запуск docker-compose(приложение, nginx и certbot)

1. Выполните клонирование репозитория:

```bash
git clone https://github.com/Srgkharkov/recordingmeet.git
```

2. Создайте файл с переменными окружения:

```bash
cd recordingmeet
echo "JWT_SECRET_KEY=your_secret_key" > .env
echo "NOTIFICATION_URL=https://example.com/notification" >> .env
```

3. Замените домен в конфигурационных файлах:

```bash
sed -i 's/example\.com/newexample.com/g' Makefile
sed -i 's/example\.com/newexample.com/g' nginx/nginx.conf
```
 
4. Выполните получение сертификатов.
```bash
make onlycertbot
```   
5. Запустите приложение.
```bash
make up
```   

### 4. Использование

#### Авторизация
Для работы API необходимо передавать валидный JWT токен в заголовке Authorization:

```bash
Authorization: Bearer <ваш-токен>
```
#### Эндпоинты
##### Начало записи конференции

Метод: GET
URL: /record?link=<ссылка_на_конференцию>
Пример: http://localhost:8080/record?link=https://meet.google.com/abc-defg-hij
Запускает запись конференции. Запись сохраняется в директории /records.
В ответе содержится идентификатор записи

##### Скачивание записей

Метод: GET
URL: /download?recordsid=<идентификатор_записи>
Пример: http://localhost:8080/download?recordsid=GM_abc-defg-hij_1234567890
Скачивает архив с файлами записей конференции.

##### Скачивание лога

Метод: GET
URL: /log?recordsid=<идентификатор_записи>
Пример: http://localhost:8080/log?recordsid=GM_abc-defg-hij_1234567890
Скачивает лог.

### 5. Тестирование
Для тестирования API можно использовать любой HTTP-клиент, например, Postman или curl.

Пример запроса для начала записи:

```bash
curl -X GET "http://localhost:8080/record?link=https://meet.google.com/abc-defg-hij" -H "Authorization: Bearer <ваш-токен>"
```

Пример запроса для скачивания архива:

```bash
curl -X GET "http://localhost:8080/download?recordsid=GM_abc-defg-hij_1234567890" -H "Authorization: Bearer <ваш-токен>" -o output.zip
```
