DRIVER=postgres
DB_NAME=postgres://localhost:5432/users_dev

migrate_up:
	goose -dir migrations $(DRIVER) $(DB_NAME) up

migrate_down:
	goose -dir migrations $(DRIVER) $(DB_NAME) down

sqlc:
	sqlc generate

swag:
	swag init --parseDependency