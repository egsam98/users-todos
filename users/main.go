package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/egsam98/users-todos/pkg/env"
	"github.com/egsam98/users-todos/users/controllers/middlewares"
	_ "github.com/egsam98/users-todos/users/docs"
	env2 "github.com/egsam98/users-todos/users/utils/env"

	_ "github.com/lib/pq"

	"github.com/egsam98/users-todos/users/controllers"
	"github.com/egsam98/users-todos/users/db"
)

// @title Users
// @version 1.0
// @BasePath /
func main() {
	var environment env2.Environment
	env.InitEnvironment(&environment)
	q := initDB(environment)
	r := initRouter(environment, q)
	log.Fatal(r.Run(environment.Addr))
}

func initRouter(environment env2.Environment, q *db.Queries) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	usersController := controllers.NewUsersController(environment, q)
	r.POST("/signup", usersController.Signup)
	r.POST("/signin", usersController.Signin)

	safeR := r.Group("/", middlewares.NewJwtMiddleware(environment, q).Process)
	safeR.GET("/users/:id", usersController.FetchUser)

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
