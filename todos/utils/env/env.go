package env

// Переменные среды
type Environment struct {
	Addr           string `env:"ADDR,required=true,default=:8081"`
	DatabaseDriver string `env:"DATABASE_DRIVER,required=true,default=postgres"`
	DatabaseConn   string `env:"DATABASE_CONN,required=true,default=postgres://localhost:5432/todos_dev"`
	AuthUrl        string `env:"AUTH_URL,required=true,default=http://localhost:8080/auth"`
}
