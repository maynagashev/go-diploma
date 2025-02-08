# Технологический стек

## Основные компоненты

### Web Framework

- **[Echo](https://github.com/labstack/echo/v4)** - быстрый и минималистичный веб-фреймворк
  - Высокая производительность
  - Простой и понятный API
  - Встроенная поддержка middleware
  - Хорошая документация

### База данных и драйверы

- **[PostgreSQL](https://www.postgresql.org/)** - основная база данных
- **[pgx](https://github.com/jackc/pgx)** - нативный PostgreSQL драйвер
  - Высокая производительность
  - Поддержка расширенных возможностей PostgreSQL
  - Встроенный пул соединений
- **[sqlx](https://github.com/jmoiron/sqlx)** - расширение стандартного database/sql
  - Поддержка сканирования в структуры
  - Типобезопасные запросы
  - Простота использования

### Конфигурация

- **[godotenv](https://github.com/joho/godotenv)** - работа с .env файлами

### Аутентификация и безопасность

- **[jwt-go v5](https://github.com/golang-jwt/jwt/v5)** - работа с JWT токенами
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - хеширование паролей

### Валидация

- **[validator v10](https://github.com/go-playground/validator)** - валидация входящих данных
  - Встроенная поддержка в Echo
  - Расширяемые правила валидации

### Логирование

- **[slog](https://pkg.go.dev/log/slog)** - структурированный логгер из стандартной библиотеки Go
  - Структурированное логирование в JSON формате
  - Поддержка уровней логирования
  - Встроенная поддержка атрибутов и контекста
  - Часть стандартной библиотеки Go 1.21+
  - Возможность кастомизации обработчиков логов

### Миграции БД

- **[goose v3](https://github.com/pressly/goose)** - управление миграциями
  - Поддержка PostgreSQL
  - Версионирование миграций
  - CLI инструмент

### Тестирование

- **[httptest](https://pkg.go.dev/net/http/httptest)** - тестирование HTTP handlers

### Линтеры и форматтеры

- **[golangci-lint](https://github.com/golangci/golangci-lint)** - набор линтеров
- **[goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)** - форматирование импортов
- **[gofmt](https://pkg.go.dev/cmd/gofmt)** - стандартный форматтер Go
