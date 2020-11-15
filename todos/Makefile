DATABASE_DRIVER ?= postgres
DATABASE_CONN_DEV ?= postgres://localhost:5432/todos_dev

migrate_up:
	goose -dir migrations $(DATABASE_DRIVER) $(DATABASE_CONN_DEV) up

migrate_down:
	goose -dir migrations $(DATABASE_DRIVER) $(DATABASE_CONN_DEV) down

sqlc:
	sqlc generate

swag:
	swag init --parseDependency

test:
	go test -v ./controllers

docker:
	docker build -f Dockerfile -t todos:latest ../