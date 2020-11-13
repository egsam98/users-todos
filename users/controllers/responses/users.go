package responses

// Ответ сервера с успешно сгенерированным токеном
type Token struct {
	Value string `json:"token"`
}

// Ответ сервера в виде найденного пользователя в БД
type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}
