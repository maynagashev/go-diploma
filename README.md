# gophermart

Накопительная система лояльности «Гофермарт». Индивидуальный дипломный проект курса «Go-разработчик»

## Быстрый старт

### 1. Подготовка окружения

```bash
# Клонирование репозитория
git clone https://github.com/your-username/gophermart.git
cd gophermart

# Создание .env файла из примера
cp .env.example .env

# Редактирование .env файла под свои настройки
nano .env
```

### 2. Запуск базы данных

```bash
# Запуск PostgreSQL через Docker
docker compose up -d pgsql

# Проверка статуса
docker compose ps
```

### 3. Сборка и запуск

```bash
# Сборка проекта
go build -o bin/gophermart cmd/gophermart/main.go

# Применение миграций
go run cmd/gophermart/main.go -m migrations

# Запуск сервера
./bin/gophermart
```

Или через флаги командной строки:
```bash
./bin/gophermart \
  -a=":8080" \
  -d="postgres://gophermart:secret@localhost:5432/gophermart?sslmode=disable" \
  -r="http://localhost:8081"
```

### 4. Проверка работоспособности

```bash
# Регистрация пользователя
curl -X POST -H "Content-Type: application/json" \
  -d '{"login":"user1","password":"pass123"}' \
  http://localhost:8080/api/user/register
```

## Конфигурация

### Переменные окружения

#### Сервер
- `RUN_ADDRESS` - адрес и порт сервера (по умолчанию `:8080`)
- `DATABASE_URI` - строка подключения к БД
- `ACCRUAL_SYSTEM_ADDRESS` - адрес системы расчета начислений

#### База данных
- `DB_DATABASE` - имя базы данных
- `DB_USERNAME` - имя пользователя БД
- `DB_PASSWORD` - пароль пользователя БД
- `FORWARD_DB_PORT` - порт для подключения к БД (по умолчанию `5432`)

### Флаги командной строки

- `-a` - адрес и порт сервера
- `-d` - строка подключения к БД
- `-r` - адрес системы расчета начислений
- `-m` - путь к директории с миграциями (по умолчанию `migrations`)

## Работа с миграциями

### Создание новой миграции

```bash
# Создание новой миграции
goose -dir migrations postgres "postgres://user:pass@localhost:5432/gophermart?sslmode=disable" create add_users_table sql
```

### Применение миграций

```bash
# Применение всех миграций
goose -dir migrations postgres "postgres://user:pass@localhost:5432/gophermart?sslmode=disable" up

# Откат последней миграции
goose -dir migrations postgres "postgres://user:pass@localhost:5432/gophermart?sslmode=disable" down

# Просмотр статуса миграций
goose -dir migrations postgres "postgres://user:pass@localhost:5432/gophermart?sslmode=disable" status
```

## Разработка

### Запуск линтеров и форматтеров

```bash
# Исправление импортов
goimports -w .

# Исправление форматирования
gofmt -w .

# Исправление длины строк
golines -w -m 120 --shorten-comments .

# Запуск линтера
golangci-lint run ./...
```

## План реализации

### 1. Структура проекта

- [x] настройка автотестов
- [x] docker compose для локальной разработки
- [x] выбор фреймворков и библиотек (echo, sqlx)
  - [x] выбор логгера (slog)
- [x] линтеры и форматтеры: `golangci-lint`, `goimports`, `gofmt`, `golines`

### 2. Основные модели данных и миграции

- [ ] `users` – пользователи
- [ ] `orders` – заказы (номера)
- [ ] `transactions` – транзакции (пополнения и списания)

### 3. Регистрация, аутентификация и авторизация пользователей

- [ ] `POST /api/user/register` — регистрация пользователя
- [ ] `POST /api/user/login` — аутентификация пользователя
- Настройка приватного ключа
- Middleware для авторизации запросов

### 4. Работа с заказами

- [ ] `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта 
  - [ ] регистрация заказа и привязка к пользователю
- [ ] `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях

### 5. Взаимодействие с системой расчета баллов лояльности

- [ ] Проверка заказа в системе accrual и начисление баллов (поллинг, воркер пул)

### 6. Баланс

- [ ] `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя

### 7. Начисление и списание баллов, получение истории списаний

- [ ] `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
- [ ] `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем

### 8. Тестирование

- [ ] юнит-тесты
- [ ] тесты API
- [ ] интеграционные тесты

### 9. Документация

- [x] `README.md` с описанием проекта и планом реализации
- [ ] API документация (swagger)

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```bash
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```bash
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Установка утилит

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/segmentio/golines@latest
# Установка pre-commit
sudo apt update
sudo apt install pipx
pipx install pre-commit
```
