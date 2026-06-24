# Тестовое задание Junior Golang Developer

REST API сервис для управления онлайн-подписками пользователей и подсчета суммарной стоимости подписок за выбранный период.

## Возможности

- Создание подписки
- Получение подписки по ID
- Получение списка подписок
- Обновление подписки
- Удаление подписки
- Подсчет суммарной стоимости подписок за выбранный период
- Фильтрация по user_id и service_name
- PostgreSQL в качестве базы данных
- SQL-миграции
- Конфигурация через переменные окружения
- Логирование HTTP-запросов и бизнес-операций
- Swagger-документация
- Запуск через Docker Compose

## Стек

- Go
- Gin
- PostgreSQL
- pgx
- golang-migrate
- slog
- swaggo
- Docker Compose

## Запуск проекта

Клонировать репозиторий:

git clone https://github.com/KinzelVA/-Junior-Golang-Developer.git

cd -Junior-Golang-Developer

Запустить сервис:

docker compose up --build

После запуска сервис будет доступен по адресам:

- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Healthcheck: http://localhost:8080/health

Запуск в фоновом режиме:

docker compose up -d --build

Остановить контейнеры:

docker compose down

Остановить контейнеры и удалить данные PostgreSQL:

docker compose down -v

## Переменные окружения

Пример конфигурации находится в файле .env.example.

APP_PORT=8080
APP_ENV=local

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSL_MODE=disable

При запуске через Docker Compose переменные окружения приложения задаются в docker-compose.yml.

## Миграции базы данных

Миграции находятся в директории migrations.

При запуске через Docker Compose миграции применяются автоматически отдельным контейнером migrate.

## Формат даты

API принимает даты подписок в формате MM-YYYY.

Пример:

start_date: 07-2025
end_date: 12-2025

В PostgreSQL даты хранятся как значения DATE. День всегда равен первому дню месяца.

Пример:

07-2025 хранится как 2025-07-01

## Логика подсчета стоимости

Стоимость подписки считается помесячно.

Пример:

Подписка стоит 400 рублей в месяц и активна с 07-2025 по 12-2025.

Выбранный период включает 6 месяцев.

Итоговая стоимость:

6 * 400 = 2400

Если end_date не указан, подписка считается активной до конца выбранного периода запроса.

## API

Базовый путь:

/api/v1

### Создать подписку

POST /api/v1/subscriptions

Пример тела запроса:

{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}

### Получить подписку по ID

GET /api/v1/subscriptions/{id}

### Получить список подписок

GET /api/v1/subscriptions

Доступные query-параметры:

- user_id
- service_name
- limit
- offset

### Обновить подписку

PUT /api/v1/subscriptions/{id}

### Удалить подписку

DELETE /api/v1/subscriptions/{id}

### Подсчитать суммарную стоимость подписок

GET /api/v1/subscriptions-total

Обязательные query-параметры:

- period_start
- period_end

Опциональные query-параметры:

- user_id
- service_name

Пример:

GET /api/v1/subscriptions-total?period_start=07-2025&period_end=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba

Пример ответа:

{
  "total": 3600
}

## Healthcheck

GET /health

Ожидаемый ответ:

{
  "database": "ok",
  "service": "subscriptions-api",
  "status": "ok"
}

## Swagger

Swagger UI доступен по адресу:

http://localhost:8080/swagger/index.html

## Проверка проекта

Запустить тесты:

go test ./...

Проверить запуск через Docker Compose:

docker compose down -v

docker compose up -d --build

docker compose ps -a

curl http://localhost:8080/health

## Структура проекта

- cmd/app — точка входа приложения
- internal/config — загрузка конфигурации
- internal/db — подключение к PostgreSQL
- internal/handler — HTTP-обработчики
- internal/model — модели, DTO и работа с датами
- internal/repository — SQL-запросы к базе данных
- internal/service — бизнес-логика
- internal/logger — настройка логирования
- migrations — SQL-миграции
- docs — Swagger-документация
