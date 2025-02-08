package main

import (
	"flag"
	"os"
	"time"
)

const (
	defaultJWTExpirationHours = 24
)

// Config содержит конфигурацию приложения.
type Config struct {
	RunAddress           string        // Адрес и порт для запуска сервера
	DatabaseURI          string        // URI базы данных
	AccrualSystemAddress string        // Адрес системы расчета начислений
	MigrationsDirectory  string        // Директория с миграциями
	JWTSecret            string        // Секретный ключ для подписи JWT токенов
	JWTExpirationPeriod  time.Duration // Период действия JWT токена
}

// parseFlags парсит флаги командной строки и переменные окружения.
func parseFlags() Config {
	var cfg Config

	// Приоритет: 1. Флаги командной строки 2. Системные переменные окружения 3. Переменные из .env
	flag.StringVar(&cfg.RunAddress, "a", getEnvOrDefault("RUN_ADDRESS", ""), "Адрес и порт для запуска сервера")
	flag.StringVar(&cfg.DatabaseURI, "d", getEnvOrDefault("DATABASE_URI", ""), "URI базы данных")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", getEnvOrDefault("ACCRUAL_SYSTEM_ADDRESS", ""), "Адрес системы расчета начислений")
	flag.StringVar(&cfg.MigrationsDirectory, "m", "migrations", "Директория с миграциями")
	flag.StringVar(&cfg.JWTSecret, "jwt-secret", getEnvOrDefault("JWT_SECRET", ""), "Секретный ключ для подписи JWT токенов")
	flag.DurationVar(
		&cfg.JWTExpirationPeriod,
		"jwt-exp",
		getDurationEnv("JWT_EXPIRATION_PERIOD", defaultJWTExpirationHours*time.Hour),
		"Период действия JWT токена",
	)

	return cfg
}

// getEnvOrDefault получает значение из переменной окружения или возвращает значение по умолчанию.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv получает значение длительности из переменной окружения.
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
