package controllers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/egsam98/users-todos/users/docs"

	"github.com/egsam98/users-todos/users/db"
	env2 "github.com/egsam98/users-todos/users/utils/env"
)

// Инициализация HTTP-маршрутов с контроллерами
func Init(environment env2.Environment, q *db.Queries) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	usersController := NewUsersController(environment, q)
	r.POST("/signup", usersController.Signup)
	r.POST("/signin", usersController.Signin)
	r.POST("/auth", usersController.Auth(false))

	safeR := r.Group("/", usersController.Auth(true))
	safeR.GET("/users/:id", usersController.FetchUser)

	return r
}
