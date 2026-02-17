.PHONY: build test lint clean install ci plugin

# Переменные
BINARY_NAME=loglinter
BIN_DIR=bin
CMD_DIR=cmd/loglinter

# Сборка
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)


plugin:
	golangci-lint custom 

# Тесты
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Запуск линтера
lint: build
	@echo "Running linter on testdata..."
	@if [ -d "testdata" ]; then \
		./$(BIN_DIR)/$(BINARY_NAME) ./testdata/src/test; \
	else \
		echo "Warning: testdata directory not found"; \
	fi

# Установка
install:
	@echo "Installing $(BINARY_NAME)..."
	go install ./$(CMD_DIR)

# Очистка
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

# CI pipeline
ci: build test lint

# Проверка кода
check:
	@echo "Running go vet..."
	go vet ./...
	@echo "Running go fmt..."
	go fmt ./...

# Все проверки
all: clean check build test lint