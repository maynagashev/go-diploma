package app

import (
	"time"
)

// Config представляет конфигурацию приложения.
type Config struct {
	DatabaseURI          string        // URI подключения к базе данных
	MigrationsDir        string        // Директория с миграциями
	RunAddress           string        // Адрес и порт для запуска сервера
	AccrualSystemAddress string        // Адрес системы расчета начислений
	JWTSecret            string        // Секретный ключ для подписи JWT токенов
	JWTExpirationPeriod  time.Duration // Период действия JWT токена
}
