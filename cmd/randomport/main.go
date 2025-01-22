package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//nolint:gosec // G102: Намеренно слушаем на всех интерфейсах для получения случайного порта
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer listener.Close()

	// Получаем адрес сокета и извлекаем номер порта
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		// Не используем log.Fatal, чтобы defer выполнился
		log.Printf("Failed to get TCP address")
		return
	}
	//nolint:forbidigo // Это утилита командной строки, использование fmt.Println допустимо
	fmt.Println(addr.Port)
}
