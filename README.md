

# Тестовое задание Junior Golang Developer

REST-сервис для управления онлайн-подписками пользователей и подсчета суммарной стоимости подписок за выбранный период.

Проект выполнен в рамках тестового задания Junior Golang Developer.

## Возможности

- CRUDL-операции над подписками:
    - создание подписки;
    - получение подписки по ID;
    - получение списка подписок;
    - обновление подписки;
    - удаление подписки.
- Подсчет суммарной стоимости подписок за выбранный период.
- Фильтрация по:
    - `user_id`;
    - `service_name`.
- PostgreSQL в качестве СУБД.
- Миграции базы данных.
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

## Формат даты

В API даты начала и окончания подписки передаются в формате:

```text
MM-YYYY

Пример:

{
  "start_date": "07-2025",
  "end_date": "12-2025"
}

В базе данные хранятся как DATE, где день всегда равен первому числу месяца.

Например:

07-2025 -> 2025-07-01
Логика подсчета стоимости

Стоимость подписки считается помесячно.

Если подписка стоит 400 рублей и активна с 07-2025 по 12-2025, то она учитывается за 6 месяцев:

6 * 400 = 2400

Если end_date не указан, подписка считается активной до конца выбранного периода запроса.

Запуск проекта
1. Клонировать репозиторий
git clone https://github.com/KinzelVA/-Junior-Golang-Developer.git
cd -Junior-Golang-Developer
2. Запустить сервис
docker compose up --build

После запуска будут доступны:

API:     http://localhost:8080
Swagger: http://localhost:8080/swagger/index.html
Health:  http://localhost:8080/health
3. Запуск в фоне
docker compose up -d --build
4. Остановка
docker compose down
5. Остановка с удалением данных PostgreSQL
docker compose down -v
Переменные окружения

Пример конфигурации находится в файле:

.env.example

Основные переменные:

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

Миграции находятся в папке:

migrations/

При запуске через Docker Compose миграции применяются автоматически отдельным контейнером migrate.

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

Пример:

GET /api/v1/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus&limit=10&offset=0
Обновить подписку
PUT /api/v1/subscriptions/{id}

Пример тела запроса:

{
  "service_name": "Yandex Plus Premium",
  "price": 600,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
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
Примеры curl-запросов
Создание подписки
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'
Получение списка
curl "http://localhost:8080/api/v1/subscriptions"
Подсчет суммы
curl "http://localhost:8080/api/v1/subscriptions-total?period_start=07-2025&period_end=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
Проверка healthcheck
curl http://localhost:8080/health

Ожидаемый ответ:

{
  "database": "ok",
  "service": "subscriptions-api",
  "status": "ok"
}
Тесты и проверка сборки
go test ./...
Swagger

Swagger UI доступен после запуска приложения:

http://localhost:8080/swagger/index.html
Структура проекта
cmd/app/                 точка входа приложения
internal/config/         загрузка конфигурации
internal/db/             подключение к PostgreSQL
internal/handler/        HTTP-обработчики
internal/model/          модели, DTO и работа с датами
internal/repository/     SQL-запросы к базе данных
internal/service/        бизнес-логика
internal/logger/         настройка логирования
migrations/              SQL-миграции
docs/                    Swagger-документация

'@ | Set-Content -Encoding UTF8 "README.md"


---

# 3. Проверяем README и проект

```powershell
go test ./...
docker compose ps -a
Invoke-RestMethod http://localhost:8080/health
git status