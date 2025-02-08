SHELL=/bin/bash
PROJECT=gophermart
PORT=8082
DB_HOST=localhost
DB_PORT=5432
DB_URI=postgres://$(PROJECT):password@$(DB_HOST):$(DB_PORT)/$(PROJECT)?sslmode=disable

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/gophermart-darwin-amd64 cmd/gophermart/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/randomport-darwin-amd64 cmd/randomport/main.go

# Локальное тестирование MacOS (Intel)
test-macos: build
	./cmd/gophermarttest/gophermarttest-darwin-amd64 \
		 -test.v -test.run=^TestGophermart$$ \
		 -gophermart-binary-path=bin/gophermart-darwin-amd64 \
		 -gophermart-host=localhost \
		 -gophermart-port=$(PORT) \
		 -gophermart-database-uri="$(DB_URI)" \
		 -accrual-binary-path=cmd/accrual/accrual_darwin_amd64 \
		 -accrual-host=localhost \
		 -accrual-port=$(shell ./bin/randomport-darwin-amd64) \
		 -accrual-database-uri="$(DB_URI)" | tee gophermarttest.log

# Локальное тестирование Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/gophermart-linux-amd64 cmd/gophermart/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/randomport-linux-amd64 cmd/randomport/main.go

test: build-linux
	./cmd/gophermarttest/gophermarttest-linux-amd64 \
		 -test.v -test.run=^TestGophermart$$ \
		 -gophermart-binary-path=bin/gophermart-linux-amd64 \
		 -gophermart-host=localhost \
		 -gophermart-port=$(PORT) \
		 -gophermart-database-uri="$(DB_URI)" \
		 -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
		 -accrual-host=localhost \
		 -accrual-port=$(shell ./bin/randomport-linux-amd64) \
		 -accrual-database-uri="$(DB_URI)" | tee gophermarttest-linux.log

perm:
	chmod -R +x bin

# Запуск сервиса
run:
	go run cmd/gophermart/main.go cmd/gophermart/flags.go cmd/gophermart/logging.go

# Запуск сервиса с локальными переменными окружения
run-env:
	RUN_ADDRESS=":$(PORT)" \
	DATABASE_URI="$(DB_URI)" \
	ACCRUAL_SYSTEM_ADDRESS="http://localhost:8081" \
	JWT_SECRET="your-256-bit-secret" \
	JWT_EXPIRATION_PERIOD="24h" \
	DEBUG=true \
	exec go run cmd/gophermart/*.go || true

# Запуск accrual сервера (blackbox)
run-accrual:
	RUN_ADDRESS=":8081" \
	DATABASE_URI="$(DB_URI)" \
	./cmd/accrual/accrual_linux_amd64

lint :
	@echo "Running linter..."
	golangci-lint run | tee lint.log