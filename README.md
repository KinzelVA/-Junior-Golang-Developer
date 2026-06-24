
# Тестовое задание Junior Golang Developer

REST-сервис для управления онлайн-подписками пользователей и подсчета суммарной стоимости подписок за выбранный период.

## Возможности

- CRUDL-операции над подписками:
    - создание подписки;
    - получение подписки по ID;
    - получение списка подписок;
    - обновление подписки;
    - удаление подписки.
- Подсчет суммарной стоимости подписок за выбранный период.
- Фильтрация подписок по `user_id` и `service_name`.
- PostgreSQL в качестве СУБД.
- SQL-миграции для инициализации базы данных.
- Конфигурация через переменные окружения.
- Логирование HTTP-запросов и бизнес-операций.
- Swagger-документация.
- Запуск через Docker Compose.

## Стек

- Go
- Gin
- PostgreSQL
- pgx
- golang-migrate
- slog
- Swagger / swaggo
- Docker Compose

## Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone https://github.com/KinzelVA/-Junior-Golang-Developer.git
cd -Junior-Golang-Developer

2. Запустить сервис
docker compose up --build

После запуска будут доступны:

API: http://localhost:8080
Swagger: http://localhost:8080/swagger/index.html
Health: http://localhost:8080/health
3. Запуск в фоне
docker compose up -d --build
4. Остановка
docker compose down
5. Остановка с удалением данных PostgreSQL
docker compose down -v
Переменные окружения

Пример конфигурации находится в файле .env.example.

APP_PORT=8080
APP_ENV=local

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSL_MODE=disable

При запуске через Docker Compose переменные окружения для приложения задаются в docker-compose.yml.

Миграции

Миграции находятся в папке migrations/.

При запуске через Docker Compose миграции применяются автоматически отдельным контейнером migrate.

Формат даты

В API даты начала и окончания подписки передаются в формате MM-YYYY.

Пример:

{
  "start_date": "07-2025",
  "end_date": "12-2025"
}

В базе данные хранятся как DATE, где день всегда равен первому числу месяца.

Пример:

07-2025 -> 2025-07-01
Логика подсчета стоимости

Стоимость подписки считается помесячно.

Если подписка стоит 400 рублей и активна с 07-2025 по 12-2025, то она учитывается за 6 месяцев:

6 * 400 = 2400

Если end_date не указан, подписка считается активной до конца выбранного периода запроса.

API

Базовый путь:

/api/v1
Создать подписку
POST /api/v1/subscriptions

Пример тела запроса:

{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
Получить подписку по ID
GET /api/v1/subscriptions/{id}
Получить список подписок
GET /api/v1/subscriptions

Доступные query-параметры:

user_id
service_name
limit
offset
Обновить подписку
PUT /api/v1/subscriptions/{id}
Удалить подписку
DELETE /api/v1/subscriptions/{id}
Подсчитать суммарную стоимость подписок
GET /api/v1/subscriptions-total

Обязательные query-параметры:

period_start
period_end

Опциональные query-параметры:

user_id
service_name

Пример:

GET /api/v1/subscriptions-total?period_start=07-2025&period_end=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba

Пример ответа:

{
  "total": 3600
}
Healthcheck
curl http://localhost:8080/health

Ожидаемый ответ:

{
  "database": "ok",
  "service": "subscriptions-api",
  "status": "ok"
}
Swagger

Swagger UI доступен после запуска приложения:

http://localhost:8080/swagger/index.html
Проверка проекта
go test ./...
Структура проекта
cmd/app/              точка входа приложения
internal/config/      загрузка конфигурации
internal/db/          подключение к PostgreSQL
internal/handler/     HTTP-обработчики
internal/model/       модели, DTO и работа с датами
internal/repository/  SQL-запросы к базе данных
internal/service/     бизнес-логика
internal/logger/      настройка логирования
migrations/           SQL-миграции
docs/                 Swagger-документация