package env

// Переменные среды
type Environment struct {
	Addr     string `env:"ADDR,required=true,default=:8081"`
	Database struct {
		Driver   string `env:"DATABASE_DRIVER,required=true,default=postgres"`
		ConnDev  string `env:"DATABASE_CONN_DEV,required=true,default=postgres://localhost:5432/todos_dev"`
		ConnTest string `env:"DATABASE_CONN_TEST,required=true,default=postgres://localhost:5432/todos_test"`
	}
	AuthUrl string `env:"AUTH_URL,required=true,default=http://localhost:8080/auth"`
}
