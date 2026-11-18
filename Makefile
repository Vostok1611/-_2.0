DB_DSN := "postgres://postgres:yourpassword@localhost:5433/gobar?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

.PHONY: migrate-new migrate migrate-down run

# создать новую пару файлов миграции: make migrate-new NAME=create_tasks
migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

# применить все миграции
migrate:
	$(MIGRATE) up

# откатить миграции
migrate-down:
	$(MIGRATE) down

# запустить приложение
run:
	go run ./cmd