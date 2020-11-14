package env

// Переменные среды
type Environment struct {
	Addr      string `env:"ADDR,required=true,default=:8080"`
	Signature string `env:"JWT_SIGNATURE,required=true,default=users-todos.io"`
	Database  struct {
		Driver   string `env:"DATABASE_DRIVER,required=true,default=postgres"`
		ConnDev  string `env:"DATABASE_CONN_DEV,required=true,default=postgres://localhost:5432/users_dev"`
		ConnTest string `env:"DATABASE_CONN_TEST,required=true,default=postgres://localhost:5432/users_test"`
	}
}
