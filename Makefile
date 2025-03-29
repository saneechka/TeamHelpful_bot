.PHONY: build run clean

# Переменные
APP_NAME=helpbot
BUILD_DIR=./build

# Сборка приложения
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/bot

# Запуск приложения
run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

# Запуск приложения без сборки
run-dev:
	@echo "Running $(APP_NAME) in dev mode..."
	@go run ./cmd/bot

# Очистка
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Тестирование
test:
	@echo "Running tests..."
	@go test -v ./... 