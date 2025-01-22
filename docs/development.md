# Руководство разработчика

## Начало работы

### 1. Настройка окружения

```bash
# Клонирование репозитория
git clone https://github.com/your-username/gophermart.git
cd gophermart

# Установка зависимостей
go mod download

# Копирование .env.example
cp .env.example .env
```

### 2. Настройка IDE

Рекомендуемые расширения для VS Code:

- Go
- EditorConfig
- GitLens

## Рабочий процесс

### Ветки

- `main` - основная ветка
- `feature/*` - новый функционал
- `fix/*` - исправления
- `docs/*` - документация

### Коммиты

Используем Conventional Commits:

```bash
feat: добавлена авторизация
fix: исправлен баг с валидацией
docs: обновлена документация API
```

## Тестирование

### Unit тесты

```bash
# Запуск всех тестов
make test

# Запуск тестов с coverage
make test-coverage
```

### Интеграционные тесты

```bash
# Запуск интеграционных тестов
make test-integration
```

## Линтинг и форматирование

```bash
# Запуск линтера
make lint

# Форматирование кода
make fmt
```

## Документация

### Swagger

```bash
# Генерация Swagger документации
make swagger
```

### Godoc

```bash
# Запуск локального сервера документации
godoc -http=:6060
```

## CI/CD

[Будет добавлено описание pipeline]

## Решение проблем

[Будет добавлен раздел с FAQ и известными проблемами]
