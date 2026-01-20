# Makefile для Pupervisor

# Переменные
BINARY_NAME=pupervisor
MAIN_PATH=cmd/server/main.go
DOCKER_IMAGE_NAME=pupervisor

# Цели
.PHONY: all build run test clean docker-build docker-run docker-stop docker-remove help

# Показать помощь
help: ## Показать это сообщение
	@echo "Доступные цели:"
	@echo "  build            - Собрать бинарный файл"
	@echo "  run              - Запустить приложение"
	@echo "  test             - Запустить тесты"
	@echo " clean            - Удалить бинарные файлы"
	@echo " docker-build     - Собрать Docker образ"
	@echo "  docker-run       - Запустить контейнер"
	@echo "  docker-stop      - Остановить контейнер"
	@echo "  docker-remove    - Удалить Docker образ"
	@echo "  help             - Показать это сообщение"

# Собрать бинарный файл
build: ## Собрать бинарный файл
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Запустить приложение
run: ## Запустить приложение
	go run $(MAIN_PATH)

# Запустить тесты
test: ## Запустить тесты
	go test ./...

# Очистить
clean: ## Удалить бинарные файлы
	rm -f $(BINARY_NAME)
	go clean

# Собрать Docker образ
docker-build: ## Собрать Docker образ
	docker build -t $(DOCKER_IMAGE_NAME) .

# Запустить Docker контейнер
docker-run: ## Запустить контейнер
	docker run -d -p 8080:8080 --name $(DOCKER_IMAGE_NAME)-container $(DOCKER_IMAGE_NAME)

# Остановить Docker контейнер
docker-stop: ## Остановить контейнер
	docker stop $(DOCKER_IMAGE_NAME)-container

# Удалить Docker образ
docker-remove: ## Удалить Docker образ
	docker rmi $(DOCKER_IMAGE_NAME)

# Установить go переменные
go-mod-tidy: ## Обновить зависимости
	go mod tidy

# Запуск с использованием .env файла
run-env: ## Запустить приложение с .env файлом
	@if [ -f .env ]; then export $$(grep -v '^#' .env | xargs) && go run $(MAIN_PATH); else echo "Файл .env не найден"; fi