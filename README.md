

# Тестовое задание Junior Golang Developer

REST-сервис для управления онлайн-подписками пользователей и подсчета суммарной стоимости подписок за выбранный период.

## Возможности

- CRUDL-операции над подписками
- Подсчет суммарной стоимости подписок за период
- Фильтрация по user_id и service_name
- PostgreSQL
- Миграции
- Логирование
- Swagger-документация
- Запуск через Docker Compose

## Запуск

```bash
docker compose up --build

После запуска:

API: http://localhost:8080
Swagger: http://localhost:8080/swagger/index.html
Health: http://localhost:8080/health
API

Базовый путь:

/api/v1
Создать подписку
POST /api/v1/subscriptions
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
Получить список подписок
GET /api/v1/subscriptions
Получить подписку по ID
GET /api/v1/subscriptions/{id}
Обновить подписку
PUT /api/v1/subscriptions/{id}
Удалить подписку
DELETE /api/v1/subscriptions/{id}
Подсчитать сумму подписок
GET /api/v1/subscriptions-total?period_start=07-2025&period_end=12-2025

Опциональные фильтры:

user_id
service_name
Формат дат

Даты передаются в формате:

MM-YYYY

Например:

07-2025

В базе дата хранится как первое число месяца:

2025-07-01
Логика подсчета

Стоимость считается помесячно.

Если подписка стоит 400 рублей и активна с 07-2025 по 12-2025, то сумма:

6 * 400 = 2400

Если end_date не указан, подписка считается активной до конца выбранного периода.

Swagger

Swagger доступен по адресу:

http://localhost:8080/swagger/index.html
Проверка
go test ./...
curl http://localhost:8080/health
Структура проекта
cmd/app/              точка входа
internal/config/      конфигурация
internal/db/          подключение к PostgreSQL
internal/handler/     HTTP handlers
internal/model/       модели и DTO
internal/repository/  SQL-запросы
internal/service/     бизнес-логика
internal/logger/      логирование
migrations/           SQL-миграции
docs/                 Swagger-документация