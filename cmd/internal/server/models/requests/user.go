package requests

// Содержит данные пользователя для авторизации
type UserRequest struct {
	Email    string
	Username string
	Password string
}
