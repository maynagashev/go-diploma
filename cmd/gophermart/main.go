package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"gophermart/internal/app"
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
	envFileLoaded := false
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
		"RUN_ADDRESS", getVarSource("RUN_ADDRESS", cfg.RunAddress, envFileLoaded),
		"DATABASE_URI", getVarSource("DATABASE_URI", cfg.DatabaseURI, envFileLoaded),
		"ACCRUAL_SYSTEM_ADDRESS", getVarSource("ACCRUAL_SYSTEM_ADDRESS", cfg.AccrualSystemAddress, envFileLoaded),
		"JWT_SECRET", maskSecret(getVarSource("JWT_SECRET", cfg.JWTSecret, envFileLoaded)),
		"JWT_EXPIRATION_PERIOD", getVarSource("JWT_EXPIRATION_PERIOD", cfg.JWTExpirationPeriod.String(), envFileLoaded),
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

	// Запускаем приложение и ждем его завершения
	if startErr := application.Start(ctx, cfg.RunAddress); startErr != nil {
		slog.Error("application error", "error", startErr)
		exitCode = 1
		return
	}

	slog.Info("application stopped")
}
