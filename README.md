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
```

### 2. Запуск базы данных

```bash
# Запуск PostgreSQL через Docker compose
docker compose up -d
```

### 3. Запуск сервиса и утилит

```bash
# Сборка и запуск сервиса gophermart
make run

# Запуск accrual
make run-accrual

# Запуск линтеров
make lint

# Запуск тестов
make test
```

## Разработка

### Установка утилит

В Linux однократно:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/segmentio/golines@latest
# Установка pre-commit
sudo apt update
sudo apt install pipx
pipx install pre-commit
```

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

- [x] `users` – пользователи
- [x] `orders` – заказы (номера)
- [ ] `transactions` – транзакции (пополнения и списания)

### 3. Регистрация, аутентификация и авторизация пользователей

- [x] `POST /api/user/register` — регистрация пользователя
- [x] `POST /api/user/login` — аутентификация пользователя
- Настройка приватного ключа
- Middleware для авторизации запросов

### 4. Работа с заказами

- [x] `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта, регистрация заказа и привязка к пользователю
- [x] `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях

### 5. Взаимодействие с системой расчета баллов лояльности

- [x] Проверка заказа в системе accrual и начисление баллов (поллинг, воркер пул)

### 6. Баланс

- [x] `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя

### 7. Начисление и списание баллов, получение истории списаний

- [x] `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
- [x] `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем

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
