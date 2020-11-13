package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"

	"github.com/egsam98/users-todos/pkg/env"
	"github.com/egsam98/users-todos/pkg/middlewares"
	"github.com/egsam98/users-todos/todos/controllers"
	"github.com/egsam98/users-todos/todos/db"
	env2 "github.com/egsam98/users-todos/todos/utils/env"
)

// @title Todos
// @version 1.0
// @BasePath /
//
// @securityDefinitions.apikey JWT-Token
// @in header
// @name Authorization
func main() {
	var environment env2.Environment
	env.InitEnvironment(&environment)
	q := initDB(environment)
	r := initRouter(environment, q)
	log.Fatal(r.Run(environment.Addr))
}

func initRouter(environment env2.Environment, q *db.Queries) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CheckAuth(environment.AuthUrl))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	todosController := controllers.NewTodosController(q)
	r.POST("/todos", todosController.CreateTodo)

	return r
}

func initDB(environment env2.Environment) *db.Queries {
	database, err := sql.Open(environment.DatabaseDriver, environment.DatabaseConn)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to database: %s\n", environment.DatabaseConn)
	return db.New(database)
}
