package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/egsam98/users-todos/users/controllers/middlewares"
	_ "github.com/egsam98/users-todos/users/docs"

	_ "github.com/lib/pq"

	"github.com/egsam98/users-todos/users/controllers"
	"github.com/egsam98/users-todos/users/db"
)

const (
	Addr         = ":8080"
	DatabaseConn = "postgres://localhost:5432/users_dev"
)

// @title Users
// @version 1.0
// @BasePath /
func main() {
	q := initDB()
	r := initRouter(q)
	log.Fatal(r.Run(Addr))
}

func initRouter(q *db.Queries) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	usersController := controllers.NewUsersController(q)
	r.POST("/signup", usersController.Signup)
	r.POST("/signin", usersController.Signin)

	safeR := r.Group("/", middlewares.NewJwtMiddleware(q).Process)
	safeR.GET("/users/:id", usersController.FetchUser)

	return r
}

func initDB() *db.Queries {
	database, err := sql.Open("postgres", DatabaseConn)
	if err != nil {
		log.Fatal(err)
	}
	return db.New(database)
}
