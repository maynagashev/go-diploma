package main

import (
	"fmt"
	"net"
)

func main() {
	// Создаем слушающий сокет на любом доступном порту
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// Закрываем слушающий сокет после завершения работы
	defer listener.Close()

	// Получаем адрес сокета и извлекаем номер порта
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Println(port)
}
