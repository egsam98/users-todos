package main

import (
	"log"

	"github.com/egsam98/users-todos/pkg/dbutils"
	"github.com/egsam98/users-todos/pkg/env"
	env2 "github.com/egsam98/users-todos/users/utils/env"

	_ "github.com/lib/pq"

	"github.com/egsam98/users-todos/users/controllers"
	"github.com/egsam98/users-todos/users/db"
)

// @title Users
// @version 1.0
// @BasePath /
//
// @securityDefinitions.apikey JWT-Token
// @in header
// @name Authorization
func main() {
	var environment env2.Environment
	env.InitEnvironment(&environment)
	q := db.New(dbutils.Init(environment.Database.Driver, environment.Database.ConnDev))
	r := controllers.Init(environment, q)
	log.Fatal(r.Run(environment.Addr))
}
