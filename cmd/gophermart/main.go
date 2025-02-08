package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"gophermart/internal/app"
)

// Глобальные переменные.
var (
	// Флаг успешной загрузки .env.
	envFileLoaded bool
)

const (
	shutdownTimeout = 10 * time.Second // Таймаут для graceful shutdown
)

func main() {
	// Код выхода по умолчанию
	exitCode := 0
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	// Загрузка .env файла, если он существует (имеет низший приоритет)
	if err := godotenv.Load(); err == nil {
		envFileLoaded = true
		slog.Debug("loaded .env file")
	}

	// Парсим флаги (имеют высший приоритет)
	cfg := parseFlags()
	flag.Parse()

	// Инициализация логгера
	initLogger()

	// Логируем все переменные окружения и их источники
	slog.Debug("configuration sources",
		"env_file_loaded", envFileLoaded,
		"RUN_ADDRESS", getVarSource("RUN_ADDRESS", cfg.RunAddress),
		"DATABASE_URI", getVarSource("DATABASE_URI", cfg.DatabaseURI),
		"ACCRUAL_SYSTEM_ADDRESS", getVarSource("ACCRUAL_SYSTEM_ADDRESS", cfg.AccrualSystemAddress),
		"JWT_SECRET", maskSecret(getVarSource("JWT_SECRET", cfg.JWTSecret)),
		"JWT_EXPIRATION_PERIOD", getVarSource("JWT_EXPIRATION_PERIOD", cfg.JWTExpirationPeriod.String()),
	)

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем канал для получения сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Запускаем горутину для обработки сигналов
	go func() {
		sig := <-sigChan
		slog.Info("received signal", "signal", sig)
		cancel() // Отменяем контекст при получении сигнала
	}()

	// Запускаем приложение
	application, err := app.New(ctx, app.Config{
		DatabaseURI:          cfg.DatabaseURI,
		MigrationsDir:        cfg.MigrationsDirectory,
		RunAddress:           cfg.RunAddress,
		AccrualSystemAddress: cfg.AccrualSystemAddress,
		JWTSecret:            cfg.JWTSecret,
		JWTExpirationPeriod:  cfg.JWTExpirationPeriod,
	})
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		exitCode = 1
		return
	}

	// Запускаем сервер в отдельной горутине
	serverErr := make(chan error, 1)
	go func() {
		if startErr := application.Start(cfg.RunAddress); startErr != nil &&
			!errors.Is(startErr, http.ErrServerClosed) {
			slog.Error("failed to start application", "error", startErr)
			serverErr <- startErr
		}
	}()

	// Ожидаем либо ошибки сервера, либо сигнала завершения
	select {
	case receivedErr := <-serverErr:
		slog.Error("server error", "error", receivedErr)
		exitCode = 1
		return
	case <-ctx.Done():
		slog.Info("shutting down server...")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if shutdownErr := application.Shutdown(shutdownCtx); shutdownErr != nil {
		slog.Error("failed to stop application", "error", shutdownErr)
		exitCode = 1
		return
	}

	slog.Info("server stopped")
}
