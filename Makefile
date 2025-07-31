      
# Имя бинарного файла, который мы собираем
BINARY_NAME=gopm

# --- Основные команды для разработки ---

build: ## Собирает Go-приложение
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

run: ## Собирает и запускает приложение с переданными аргументами (Пример: make run ARGS="create ...")
	@make build
	@echo "Running $(BINARY_NAME) with args: $(ARGS)"
	@./$(BINARY_NAME) $(ARGS)

test: ## Запускает все тесты в проекте
	@echo "Running tests..."
	@go test ./... -v

clean: ## Удаляет скомпилированный бинарник
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

tidy: ## Приводит в порядок зависимости в go.mod
	@echo "Tidying go modules..."
	@go mod tidy


# --- Команды для работы с тестовым SFTP-сервером ---

sftp-up: ## Запускает тестовый SFTP-сервер в Docker
	@echo "Starting SFTP server container..."
	@docker rm -f sftp-server > /dev/null 2>&1 || true
	@docker run -p 2222:22 -d --name sftp-server atmoz/sftp testuser:testpass:::upload
	@echo "SFTP server is running on localhost:2222 (user: testuser, pass: testpass)"

sftp-down: ## Останавливает и удаляет тестовый SFTP-сервер
	@echo "Stopping and removing SFTP server container..."
	@docker stop sftp-server > /dev/null 2>&1
	@docker rm sftp-server > /dev/null 2>&1

sftp-ls: ## Показывает список файлов на тестовом SFTP-сервере
	@echo "Files on SFTP server:"
	@docker exec sftp-server ls -l /home/testuser/upload


# --- Команда-справка ---

help: ## Показывает эту справку
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build run test clean tidy sftp-up sftp-down sftp-ls help

    