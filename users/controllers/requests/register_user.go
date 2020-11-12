package requests

// Тело запроса для регистрации нового пользователя в системе
type RegisterUser struct {
	Username             string `json:"username" binding:"required,max=12" error:"must be non empty and have max length 12"`
	Password             string `json:"password" binding:"required,min=6,max=12" error:"must have length 6..12 symbols"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password" error:"must match password field"`
}
