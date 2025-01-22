package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort       = "8080"
	defaultTimeout    = 30 * time.Second
	readHeaderTimeout = 2 * time.Second
)

func main() {
	initLogger()
	port := defaultPort

	// Создаем новый мультиплексор для маршрутизации
	mux := http.NewServeMux()

	// Обработчик запросов
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		// отправка ответа
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello, Gopher!")); err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	})

	// Создаем HTTP сервер с таймаутами
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Старт сервера с выводом информации
	slog.Info("Starting server", "port", port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}

func initLogger() {
	// Создаем переменную для уровня логирования и устанавливаем ее в Info
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	// Создаем новый обработчик с настроенным уровнем логирования
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Устанавливаем созданный логгер как логгер по умолчанию
	slog.SetDefault(logger)
}
