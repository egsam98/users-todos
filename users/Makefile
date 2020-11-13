DATABASE_DRIVER ?= postgres
DATABASE_CONN ?= postgres://localhost:5432/users_dev

migrate_up:
	goose -dir migrations $(DATABASE_DRIVER) $(DATABASE_CONN) up

migrate_down:
	goose -dir migrations $(DATABASE_DRIVER) $(DATABASE_CONN) down

sqlc:
	sqlc generate

swag:
	swag init --parseDependency