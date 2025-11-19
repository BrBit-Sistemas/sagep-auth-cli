.PHONY: build install test clean help

# Nome do binário
BINARY_NAME=sagep-auth-cli
CMD_PATH=./cmd/sagep-auth-cli

help: ## Exibe esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compila o CLI
	@echo "Compilando $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build concluído: ./$(BINARY_NAME)"

install: ## Instala o CLI globalmente
	@echo "Instalando $(BINARY_NAME)..."
	@go install $(CMD_PATH)
	@echo "Instalado com sucesso!"

test: ## Executa os testes
	@echo "Executando testes..."
	@go test ./...

clean: ## Remove arquivos gerados
	@echo "Limpando arquivos..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Limpeza concluída!"

run: ## Executa o CLI localmente (exemplo: make run ARGS="sync --manifest ./auth-manifest.example.yaml")
	@go run $(CMD_PATH) $(ARGS)

