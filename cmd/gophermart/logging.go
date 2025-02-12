package main

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	minSecretLength = 4 // Минимальная длина секрета для маскирования
)

// initLogger инициализирует логгер.
func initLogger() {
	// Настраиваем уровень логирования
	var programLevel = new(slog.LevelVar)
	if os.Getenv("DEBUG") != "" {
		programLevel.Set(slog.LevelDebug)
	}

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	logger := slog.New(h).With(
		"environment", os.Getenv("ENV"),
	)
	slog.SetDefault(logger)
}

// getVarSource возвращает значение переменной и её источник.
func getVarSource(name string, value string, envFileLoaded bool) string {
	// Проверяем наличие переменной в окружении
	if envValue := os.Getenv(name); envValue != "" {
		if envFileLoaded {
			return fmt.Sprintf("%s (from .env)", value)
		}
		return fmt.Sprintf("%s (from environment)", value)
	}
	// Если значение есть, но его нет в окружении, значит оно из флага
	if value != "" {
		return fmt.Sprintf("%s (from flag)", value)
	}
	return "not set"
}

// maskSecret маскирует секретные значения для логов.
func maskSecret(s string) string {
	if len(s) <= minSecretLength {
		return "***"
	}
	return s[:2] + "***" + s[len(s)-2:]
}
