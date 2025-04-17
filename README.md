# auth-test-task
Тестовое задание от Medods по созданию части сервиса аутентификации.
## Использованные технологии и библиотеки.
Технологии: Docker, Docker-compose, Jwt, PostgreSQL.

Библиотеки: Chi в качестве роутера, Zap для логгирования, Squirrel для составления sql-запросов
## Установка проекта и скачивание зависимостей
```bash
   git clone github.com/solumD/auth-test-task
   cd auth-test-task/
   make install-deps
   go mod tidy
```

## Запуск (docker-compose обязателен)
Поменять значения в .env-файле, если необходимо.
```dotenv
  PG_DATABASE_NAME=auth
  PG_USER=auth-user
  PG_PASSWORD=auth-password
  PG_PORT=54321
  MIGRATION_DIR=./migrations
  
  PG_DSN="host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
  MIGRATION_DSN="host=pg port=5432 dbname=auth user=auth-user password=auth-password sslmode=disable"
  
  LOGGER_LEVEL=info #info, error, debug
  
  SERVER_HOST=localhost
  SERVER_PORT=8080
  
  JWT_KEY="gO-Is-AwEsOmE"
```

Выполнить в терминале:
```bash
  docker compose up -d
  go run cmd/app/main.go
```
БД Postgres поднимается в отдельном docker-контейнере. При запуске мигратор автоматически накатывает миграции.

## Эндпоинты
#### GET /token/generate - получение access и refresh токенов по guid пользователя.
 
##### Example Input (query param): 
```
  ?guid=03e50099-683f-49f1-9ac2-7a04d7a8a105
```

##### Example Response (Ok): 
```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4OTc3OTksIlVzZXJHVUlEIjoiZWUxNjRhZDQtNTA5OC00NDM5LWI0ZjYtNzA5YjUzZWU0NjA5IiwiVXNlcklQIjoiMTI3LjAuMC4xIiwiQWNjZXNzVG9rZW5VSUQiOiI5MGMzYzhlNC01NmZkLTQxY2YtOWFiMy03Yzg5ODExMzBiNTIifQ.ZTfEuiK3fukQk_P0vxSlEWdiqp1WA4uukP2Di_xEw46ZVgKRIqeqg2PbPkAzLCkmmP7dC16_6UwSDHh8gcZ6zg",
    "refresh_token": "OTQ1NDM4ZTYtYjYwNi00ZDBkLTlmMzgtYjY4NjRhOTMzZDc4"
}
```

##### Example Response (Error): 
```
{
     "error_message": "failed to get user ip"
}
```

#### POST /token/refresh - обновление access и refresh токенов

##### Example Request: 
```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4OTc1NzUsIlVzZXJHVUlEIjoiZWUxNjRhZDQtNTA5OC00NDM5LWI0ZjYtNzA5YjUzZWU0NjA5IiwiVXNlcklQIjoiMTI3LjAuMC4xIiwiQWNjZXNzVG9rZW5VSUQiOiJmZmM3MDA3Ny1jM2I4LTRiY2EtYTg5ZS1hZmUwNGE5NTExNjgifQ.ylT6fACIxaQ2Lmp0hem6Nc-BEPyAzhMHvd9iNbov1B-7uEQ9OcZH8lQV-emj04xH3GfhwN58W3HIBvvAsjIxIg",
    "refresh_token": "OGM5MTI0YWMtODMzYi00MTdiLTkxZmMtMzY3MjIwZGRiOGJm"
}
```

##### Example Response (OK): 
```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4OTU4NzYsIlVzZXJHVUlEIjoiZWUxNjRhZDQtNTA5OC00NDM5LWI0ZjYtNzA5YjUzZWU0NjA5IiwiVXNlcklQIjoiMTI3LjAuMC4xIiwiQWNjZXNzVG9rZW5VSUQiOiI1ZTVjZjE5My0yY2UwLTQyMDItYWZiOS0wYTY4Y2U0YzllZTkifQ.3aHvQ4BRT1pf8Nz1hEaADG-fsRLt93umTgU9ZFgVZJG9huAsW-1zeieYIicXhkLe4tdjTsONfqEeuGXrniwPUA",
    "refresh_token": "NmVjN2M5NDctMjM4Yi00MTA4LTk5ZTAtYjE5NjYyNmUwNWY3"
}
```

##### Example Response (Error): 
```
{
    "error_message": "old and curr refresh tokens's do not match"
}
```
