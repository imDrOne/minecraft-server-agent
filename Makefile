export
	GO111MODULE=on
	APP_CONFIG_DIR=./configs

BINARY_NAME=minecraft-agent
MAIN_PATH=./cmd/main.go
BUILD_DIR=./build
GO_FILES=$(shell find . -name "*.go" -type f)

# Версия и информация о сборке
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME = $(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

.PHONY: help build run test clean deps lint fmt vet install docker-build

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# -- ✴️ Core ---
deps: ## Установить зависимости
	@echo "Установка зависимостей..."
	go mod download
	go mod tidy

build: deps ## Собрать бинарный файл
	@echo "Сборка проекта..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Сборка завершена: $(BUILD_DIR)/$(BINARY_NAME)"

build-dev: ## Быстрая сборка для разработки
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: ## Запустить приложение
	go run $(MAIN_PATH)

clean: ## Очистить сгенерированные файлы
	@echo "Очистка..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	go clean

# --- 🧪 Tests ---
test: ## Запустить тесты
	@echo "Запуск тестов..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Запустить тесты с покрытием
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет о покрытии: coverage.html"

bench: ## Запустить бенчмарки
	go test -bench=. -benchmem ./...

# --- 💅 Linting & etc. ---
lint: ## Запустить golangci-lint
	golangci-lint run ./...

fmt: ## Форматировать код
	go fmt ./...
	goimports -w .

vet: ## Запустить go vet
	go vet ./...

check: fmt vet lint test ## Полная проверка кода

# 🐳 Docker
docker-build: ## Собрать Docker образ
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run: docker-build ## Запустить в Docker контейнере
	docker run --rm -it $(BINARY_NAME):$(VERSION)
