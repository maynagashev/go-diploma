package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"gophermart/internal/app"
)

var (
	runAddress          string
	databaseURI         string
	accrualSystemAddr   string
	migrationsDirectory string
	jwtSecret           string
	jwtExpirationPeriod time.Duration
)

func init() {
	// Загрузка .env файла, если он существует
	_ = godotenv.Load()

	// Флаги командной строки
	flag.StringVar(&runAddress, "a", os.Getenv("RUN_ADDRESS"), "Адрес и порт для запуска сервера")
	flag.StringVar(&databaseURI, "d", os.Getenv("DATABASE_URI"), "URI базы данных")
	flag.StringVar(&accrualSystemAddr, "r", os.Getenv("ACCRUAL_SYSTEM_ADDRESS"), "Адрес системы расчета начислений")
	flag.StringVar(&migrationsDirectory, "m", "migrations", "Директория с миграциями")
	flag.StringVar(&jwtSecret, "jwt-secret", os.Getenv("JWT_SECRET"), "Секретный ключ для подписи JWT токенов")
	flag.DurationVar(
		&jwtExpirationPeriod,
		"jwt-exp",
		getDurationEnv("JWT_EXPIRATION_PERIOD", 24*time.Hour),
		"Период действия JWT токена",
	)
}

// getDurationEnv получает значение длительности из переменной окружения
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func main() {
	flag.Parse()

	// Установка значений по умолчанию
	if runAddress == "" {
		runAddress = ":8080"
	}
	if databaseURI == "" {
		slog.Error("требуется указать URI базы данных")
		os.Exit(1)
	}
	if accrualSystemAddr == "" {
		slog.Error("требуется указать адрес системы расчета начислений")
		os.Exit(1)
	}
	if jwtSecret == "" {
		slog.Error("требуется указать секретный ключ для JWT")
		os.Exit(1)
	}

	// Создание контекста, который слушает сигналы прерывания от ОС
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инициализация приложения
	application, err := app.New(ctx, app.Config{
		DatabaseURI:          databaseURI,
		MigrationsDir:        migrationsDirectory,
		RunAddress:           runAddress,
		AccrualSystemAddress: accrualSystemAddr,
		JWTSecret:            jwtSecret,
		JWTExpirationPeriod:  jwtExpirationPeriod,
	})
	if err != nil {
		slog.Error("не удалось инициализировать приложение", "error", err)
		os.Exit(1)
	}

	// Запуск приложения
	go func() {
		if err := application.Start(runAddress); err != nil {
			slog.Error("не удалось запустить сервер", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала прерывания
	<-ctx.Done()

	// Корректное завершение работы
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		slog.Error("не удалось корректно завершить работу сервера", "error", err)
		os.Exit(1)
	}
}

func initLogger() {
	// Создаем переменную для уровня логирования и устанавливаем ее в Debug
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	// Создаем новый обработчик с настроенным уровнем логирования
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Устанавливаем созданный логгер как логгер по умолчанию
	slog.SetDefault(logger)
}
