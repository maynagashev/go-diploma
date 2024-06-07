SHELL=/bin/bash
PROJECT=gophermart
DB_HOST=localhost
DB_PORT=5432
DB_URI=postgres://$(PROJECT):password@$(DB_HOST):$(DB_PORT)/$(PROJECT)?sslmode=disable

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/gophermart-darwin-amd64 cmd/gophermart/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/randomport-darwin-amd64 cmd/randomport/main.go

# Локальное тестирование MacOS (Intel)
autotests-darwin-amd64: build
	./cmd/gophermarttest/gophermarttest-darwin-amd64 \
		 -test.v -test.run=^TestGophermart$$ \
		 -gophermart-binary-path=bin/gophermart-darwin-amd64 \
		 -gophermart-host=localhost \
		 -gophermart-port=8080 \
		 -gophermart-database-uri="$(DB_URI)" \
		 -accrual-binary-path=cmd/accrual/accrual_darwin_amd64 \
		 -accrual-host=localhost \
		 -accrual-port=$(shell ./bin/randomport-darwin-amd64) \
		 -accrual-database-uri="$(DB_URI)" | tee gophermarttest.log

perm:
	chmod -R +x bin