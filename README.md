# AIQADAM Backend

Go API для платформы AIQADAM: авторизация (JWT + bcrypt), админка, курсы, кабинет сотрудника, файловое хранилище.

## Стек

- **Go** + [Echo](https://echo.labstack.com/)
- **PostgreSQL** — [pgx](https://github.com/jackc/pgx)
- **Миграции** — [golang-migrate](https://github.com/golang-migrate/migrate)
- **Docker** — multi-stage build + Compose (Postgres + API)
- **Конфиг** — переменные окружения + [godotenv](https://github.com/joho/godotenv) для локальной разработки

## Структура (clean architecture)

```
backend/
  cmd/api/              # HTTP API
  cmd/migrate/          # CLI миграций (-command up|down|status)
  internal/
    config/             # загрузка env
    database/           # PostgreSQL pool
    migrate/            # обёртка golang-migrate
    http/               # Echo server, routes, middleware
    auth/, users/, courses/
  migrations/           # SQL-миграции
  scripts/import/       # Seed-данные и импорт из дампа
  docker-compose.yml
  Dockerfile
  docker-entrypoint.sh
```

## Docker запуск

Из каталога `backend`:

```sh
docker compose up --build
```

Или в фоне:

```sh
docker compose up --build -d
```

Через Makefile:

```sh
make compose-up
make compose-up-d
```

### Проверка

```sh
curl http://localhost:8080/
curl http://localhost:8080/healthz
```

Ожидаемый healthcheck (БД доступна):

```json
{
  "status": "ok",
  "database": "ok"
}
```

### Логи

```sh
docker compose logs -f backend
```

В логах при старте backend-контейнера:

- `migrations started`
- `migrations completed` (или `migrations applied successfully`)
- `api server started`

### Остановка

```sh
docker compose down
```

Полная очистка базы и volumes:

```sh
docker compose down -v
```

### Сеть и порты

| Сервис | С хоста | Внутри Docker network |
|--------|---------|------------------------|
| PostgreSQL | `localhost:54329` | `postgres:5432` |
| Backend API | `localhost:8080` | `backend:8080` |

Backend внутри Compose подключается к БД по hostname **`postgres`**, не `localhost`.

`DATABASE_URL` в `docker-compose.yml`:

```
postgres://postgres:postgres@postgres:5432/aiqadam?sslmode=disable
```

Миграции применяются **автоматически** при старте backend-контейнера (`docker-entrypoint.sh` → `migrate up` → `api`).

## Локальный запуск (без Docker)

1. Поднимите только Postgres (опционально):

```sh
docker compose up postgres -d
```

2. Скопируйте env:

```sh
cp .env.example .env
```

3. Для `go run` используйте `localhost:54329`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:54329/aiqadam?sslmode=disable
```

4. Миграции и API:

```sh
go mod tidy
make migrate-up
make run
```

Сервер: **http://localhost:8080**

## Makefile

| Команда | Описание |
|---------|----------|
| `make tidy` | `go mod tidy` |
| `make run` | Локальный API (`go run ./cmd/api`) |
| `make migrate-up` | Миграции вверх (читает `.env`) |
| `make migrate-down` | Откат одной миграции |
| `make migrate-status` | Текущая версия миграций |
| `make compose-up` | `docker compose up --build` |
| `make compose-up-d` | Compose в фоне |
| `make compose-down` | Остановить контейнеры |
| `make compose-down-v` | Остановить + удалить volumes |
| `make compose-logs` | Логи в follow-режиме |
| `make compose-ps` | Статус контейнеров |
| `make import-data` | Импорт `scripts/import/local_data.sql` |

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `PORT` | HTTP-порт | `8080` |
| `APP_ENV` | Окружение | `development` |
| `DATABASE_URL` | PostgreSQL connection string | обязательно |
| `MIGRATIONS_PATH` | Путь к SQL-миграциям | `migrations` или `/app/migrations` |
| `AUTH_JWT_SECRET` | Секрет JWT (будущее) | — |
| `AUTH_JWT_ISSUER` | Issuer JWT | — |
| `AUTH_JWT_EXPIRY_HOURS` | Срок жизни JWT (часы) | `24` |
| `CORS_ALLOWED_ORIGINS` | Origins через запятую | `*` если пусто |

## Auth (вместо Supabase Auth)

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| `POST` | `/auth/login` | — | Email + password → JWT |
| `GET` | `/functions/v1/auth/me` | Bearer JWT | Текущий пользователь + роль из `profiles` |
| `GET` | `/functions/v1/admin/users` | Admin JWT | Список пользователей (`?search=`) |
| `GET` | `/functions/v1/admin/users/:id` | Admin JWT | Профиль пользователя |
| `POST` | `/functions/v1/admin-add-user` | Admin JWT | Создать пользователя |
| `POST` | `/functions/v1/admin-update-user` | Admin JWT | Обновить профиль / пароль |
| `POST` | `/functions/v1/admin-delete-user` | Admin JWT | Удалить пользователей |

### POST /auth/login

```json
{ "email": "adminaq@gmail.com", "password": "your-password" }
```

Ответ:

```json
{
  "access_token": "<jwt>",
  "token_type": "bearer",
  "expires_in": 86400,
  "user": { "id": "...", "email": "...", "role": "super_admin" }
}
```

Пароли из Supabase (`auth.users.encrypted_password`) — bcrypt, совместимы с `bcrypt.Compare`.

### Импорт seed-данных

После `docker compose up` и миграций (из каталога `backend/`):

```powershell
make import-data
# или
.\scripts\import\import-local-data.ps1
```

Подробнее: [scripts/import/README.md](scripts/import/README.md).

### CORS и порт 8080

В `APP_ENV=development` разрешены все `http://localhost:*` и `http://127.0.0.1:*`.

Если Docker не стартует из‑за занятого **8080** — остановите локальный `go run ./cmd/api` или задайте в `.env`:

```env
BACKEND_HOST_PORT=8081
PUBLIC_API_URL=http://localhost:8081
```

## Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/` | Статус backend |
| `GET` | `/healthz` | Healthcheck + проверка БД (`503` если БД недоступна) |
| `POST` | `/auth/login` | Вход, выдача JWT |
| `GET` | `/functions/v1/auth/me` | Профиль по JWT |

### Примеры ответов

**GET /** — `200`

```json
{ "name": "AIQADAM Backend", "status": "running" }
```

**GET /healthz** — `200` если БД доступна

```json
{ "status": "ok", "database": "ok" }
```

**GET /healthz** — `503` если БД недоступна

```json
{ "status": "error", "database": "error" }
```

**GET /functions/v1/auth/me** — `200`

```json
{
  "id": "...",
  "email": "adminaq@gmail.com",
  "role": "super_admin",
  "is_active": true
}
```

## Архитектура

```
Frontend (../frontend) → Backend (Go/Echo) → PostgreSQL
```
