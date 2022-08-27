package requests // TODO: переместить в models

// Содержит данные пользователя для авторизации
type UserRequest struct {
	Email    string
	Username string
	Password string
}
