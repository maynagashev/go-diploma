package main

import "net/http"

func main() {

	// обработчик запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// отправка ответа
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Gopher!"))
	})

	// старт сервера на порту 8080 без маршрутов и ответом 200 OK
	http.ListenAndServe(":8080", nil)
}
