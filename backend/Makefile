.PHONY: run build test lint migrate-up migrate-down migrate-create clean

# Carregar variáveis do .env se existir
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Caminho do GOPATH/bin
GOPATH_BIN := $(HOME)/go/bin
MIGRATE := $(GOPATH_BIN)/migrate

# Executar a aplicação
run:
	go run cmd/server/main.go

# Build da aplicação
build:
	go build -o bin/server cmd/server/main.go

# Executar testes
test:
	go test -v ./...

# Executar testes com cobertura
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint do código
lint:
	golangci-lint run

# Executar migrations (up)
migrate-up:
	$(MIGRATE) -path ./migrations -database "postgres://$(DB_USER_DEV):$(DB_PASSWORD_DEV)@$(DB_HOST_DEV):$(DB_PORT_DEV)/$(DB_NAME_DEV)?sslmode=disable" up

# Executar migrations (down)
migrate-down:
	$(MIGRATE) -path ./migrations -database "postgres://$(DB_USER_DEV):$(DB_PASSWORD_DEV)@$(DB_HOST_DEV):$(DB_PORT_DEV)/$(DB_NAME_DEV)?sslmode=disable" down

# Criar nova migration
migrate-create:
	@read -p "Nome da migration: " name; \
	$(MIGRATE) create -ext sql -dir ./migrations -seq $$name

# Limpar arquivos gerados
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Instalar dependências
deps:
	go mod download
	go mod tidy

# Instalar ferramentas de desenvolvimento
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Gerar documentação Swagger
swagger:
	swag init -g cmd/server/main.go -o ./docs