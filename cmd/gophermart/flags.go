package main

import (
	"flag"
	"os"
	"time"
)

const (
	defaultJWTExpirationHours = 24
)

var (
	// Флаги командной строки.
	runAddress          string        // Адрес и порт для запуска сервера
	databaseURI         string        // URI базы данных
	accrualSystemAddr   string        // Адрес системы расчета начислений
	migrationsDirectory string        // Директория с миграциями
	jwtSecret           string        // Секретный ключ для подписи JWT токенов
	jwtExpirationPeriod time.Duration // Период действия JWT токена
)

func initFlags() {
	flag.StringVar(&runAddress, "a", os.Getenv("RUN_ADDRESS"), "Адрес и порт для запуска сервера")
	flag.StringVar(&databaseURI, "d", os.Getenv("DATABASE_URI"), "URI базы данных")
	flag.StringVar(&accrualSystemAddr, "r", os.Getenv("ACCRUAL_SYSTEM_ADDRESS"), "Адрес системы расчета начислений")
	flag.StringVar(&migrationsDirectory, "m", "migrations", "Директория с миграциями")
	flag.StringVar(&jwtSecret, "jwt-secret", os.Getenv("JWT_SECRET"), "Секретный ключ для подписи JWT токенов")
	flag.DurationVar(
		&jwtExpirationPeriod,
		"jwt-exp",
		getDurationEnv("JWT_EXPIRATION_PERIOD", defaultJWTExpirationHours*time.Hour),
		"Период действия JWT токена",
	)
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
