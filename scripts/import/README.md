# Импорт данных в локальный Postgres

Запускайте команды из каталога **`backend/`**.

## Быстрый старт

```powershell
cd backend
docker compose up -d --build
.\scripts\import\import-local-data.ps1 -SkipMigrate
# или после миграций контейнером:
make import-data
```

## Файлы

| Файл | Описание |
|------|----------|
| `local_data.sql` | Готовый seed (пользователи, курсы, уроки) |
| `data_only_dump.sql` | Полный дамп из Supabase (опционально) |
| `import-local-data.ps1` | Миграции + импорт |
| `extract-from-dump.ps1` | Пересобрать `local_data.sql` из дампа |

## Тестовый вход

| Email | Роль |
|-------|------|
| `adminaq@gmail.com` | super_admin |
| `u@gmail.com` | user |

Пароли — те же, что были в Supabase (bcrypt в `auth.users`).
