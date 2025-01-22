# Развертывание

## Требования

### Системные требования

- Go 1.23+
- PostgreSQL 14+
- Docker для запуска БД
- Make для сборки и запуска бинарников, тестов, миграций

### Переменные окружения

```env
# Сервер
RUN_ADDRESS=:8080
DATABASE_URI=postgres://user:pass@localhost:5432/dbname
ACCRUAL_SYSTEM_ADDRESS=http://localhost:8081

# Настройки для развертывания сервиса локально в docker-compose
DB_DATABASE=gophermart
DB_USERNAME=gophermart
DB_PASSWORD=password
FORWARD_DB_PORT=5432 
```

## Локальное развертывание

### 1. Через Docker Compose

```bash
# Запуск БД в докере
docker compose up -d

# Остановка БД в докере
docker compose down
```

### 2. Бинарники

```bash
# Установка зависимостей
go mod download

# Применение миграций
make migrate

# Запуск сервера
make run
```

## Миграции

[Будет добавлено позже]
