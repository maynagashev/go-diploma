package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	initLogger()
	port := "8080"

	// Обработчик запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// отправка ответа
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Gopher!"))
	})

	// Старт сервера с выводом информации
	slog.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
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
