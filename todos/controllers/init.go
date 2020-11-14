package controllers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/egsam98/users-todos/todos/docs"

	"github.com/egsam98/users-todos/pkg/middlewares"
	"github.com/egsam98/users-todos/todos/db"
	env2 "github.com/egsam98/users-todos/todos/utils/env"
)

// Инициализация маршрутов и контролеров
func Init(environment env2.Environment, q *db.Queries) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	safeR := r.Group("/", middlewares.CheckAuth(environment.AuthUrl))

	todosController := NewTodosController(q)
	safeR.POST("/todos", todosController.CreateTodo)
	safeR.PUT("/todos/:id", todosController.UpdateTodo)
	safeR.DELETE("/todos/:id", todosController.DeleteTodo)
	safeR.GET("/todos", todosController.FetchAll)
	safeR.POST("/todos/before", todosController.FetchBeforeDeadline)

	return r
}
