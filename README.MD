## Для запуска:
```shell script
docker-compose up
```

## Потребуются дополнительно:
* Для локальной работы с миграциями
```shell script
go get -u github.com/pressly/goose/cmd/goose
```
* Для генерации Swagger-документации 2.0:
```shell script
go get -u github.com/swaggo/swag/cmd/swag
```
* Для генерации sql запросов в виде go кода
https://github.com/kyleconroy/sqlc

#### macOS

```
brew install kyleconroy/sqlc/sqlc
```

#### Ubuntu

```
sudo snap install sqlc
```

#### go get

```
go get github.com/kyleconroy/sqlc/cmd/sqlc
```

#### Docker

```
docker pull kjconroy/sqlc
```

Run `sqlc` using `docker run`:

```
docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
```
*