APP_NAME := learn-gin

DB_DRIVER := mysql
DB_DSN := root:root@tcp(127.0.0.1:3306)/phase4_gin?charset=utf8mb4&parseTime=True&loc=Local
MIGRATIONS_DIR := migrations

# 默认目标
.PHONY: help
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Migration:"
	@echo "	 make migrate-up		Run database migrations"
	@echo "  make migrate-down 		Rollback last migration"
	@echo "  make migrate-status 	Show migration status"
	@echo ""
	@echo "Server:"
	@echo "  make run 				Run HTTP server"
	@echo ""
	@echo "Tools:"
	@echo "  make tidy              go mod tidy"
	@echo ""

.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" down

.PHONY: migrate-status
migrate-status:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" status

.PHONY: run
run:
	go run cmd/main.go

.PHONY: tidy
tidy:
	go mod tidy