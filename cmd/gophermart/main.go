package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"

	// обработчик запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// отправка ответа
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Gopher!"))
	})

	// старт сервера с выводом информации
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
