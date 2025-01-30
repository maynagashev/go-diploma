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

// Глобальные переменные
var (
	// Флаг успешной загрузки .env
	envFileLoaded bool
)

func init() {
	// Загрузка .env файла, если он существует
	if err := godotenv.Load(); err == nil {
		envFileLoaded = true
	}

	// Инициализация флагов
	initFlags()
}

func main() {
	// Парсим флаги
	flag.Parse()

	// Инициализация логгера
	initLogger()

	// Логируем все переменные окружения и их источники
	slog.Debug("configuration sources",
		"env_file_loaded", envFileLoaded,
		"RUN_ADDRESS", getVarSource("RUN_ADDRESS", runAddress),
		"DATABASE_URI", getVarSource("DATABASE_URI", databaseURI),
		"ACCRUAL_SYSTEM_ADDRESS", getVarSource("ACCRUAL_SYSTEM_ADDRESS", accrualSystemAddr),
		"JWT_SECRET", maskSecret(getVarSource("JWT_SECRET", jwtSecret)),
		"JWT_EXPIRATION_PERIOD", getVarSource("JWT_EXPIRATION_PERIOD", jwtExpirationPeriod.String()),
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
		DatabaseURI:          databaseURI,
		MigrationsDir:        migrationsDirectory,
		RunAddress:           runAddress,
		AccrualSystemAddress: accrualSystemAddr,
		JWTSecret:            jwtSecret,
		JWTExpirationPeriod:  jwtExpirationPeriod,
	})
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Запускаем сервер в отдельной горутине
	serverErr := make(chan error, 1)
	go func() {
		if err := application.Start(runAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start application", "error", err)
			serverErr <- err
		}
	}()

	// Ожидаем либо ошибки сервера, либо сигнала завершения
	select {
	case err := <-serverErr:
		slog.Error("server error", "error", err)
		os.Exit(1)
	case <-ctx.Done():
		slog.Info("shutting down server...")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to stop application", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
