package dtos

// Модель токена для работы с БД
// Содержит ссылку на пользователя и refresh_token
type TokenDto struct {
	User         string
	RefreshToken string
}
