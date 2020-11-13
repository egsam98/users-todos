package responses

// Ответ сервера с успешно сгенерированным токеном
type Token struct {
	Value string `json:"token"`
}
