package env

// Переменные среды
type Environment struct {
	Addr           string `env:"ADDR,required=true,default=:8080"`
	DatabaseDriver string `env:"DATABASE_DRIVER,required=true,default=postgres"`
	DatabaseConn   string `env:"DATABASE_CONN,required=true,default=postgres://localhost:5432/users_dev"`
	Signature      string `env:"JWT_SIGNATURE,required=true,default=users-todos.io"`
}
