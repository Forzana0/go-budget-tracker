# Підтягуємо змінні з файлу .env, якщо він існує
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# .PHONY вказує make, що це назви команд, а не реальні файли на диску
.PHONY: help run build migrate-up migrate-down migrate-version migrate-create

# Команда за замовчуванням (показує підказку)
help:
	@echo "Доступні команди:"
	@echo "  make run             - Запустити додаток"
	@echo "  make build           - Скомпілювати бінарний файл"
	@echo "  make migrate-up      - Накатити всі нові міграції"
	@echo "  make migrate-down    - Відкотити останню міграцію (1 крок)"
	@echo "  make migrate-version - Показати поточну версію схеми БД"
	@echo "  make migrate-create name=<назва> - Створити нову пару міграцій"

# Запуск Go сервера
run:
	go run main.go

# Збірка бінарного файлу в папку bin/
build:
	go build -o bin/budget-tracker main.go

# Накатити всі доступні міграції
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

# Відкотити рівно одну останню міграцію
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

# Перевірити поточний номер версії бази даних
migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version

# Створити нову пару міграцій (наприклад: make migrate-create name=add_user_email)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Помилка: вкажіть назву міграції. Приклад: make migrate-create name=add_user_email"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)
	