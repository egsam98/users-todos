package dbutils

import (
	"database/sql"
	"log"
)

// Инициализация сессии к БД
func Init(driver, conn string) *sql.DB {
	database, err := sql.Open(driver, conn)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to database: %s\n", conn)
	return database
}
