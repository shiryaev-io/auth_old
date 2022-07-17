package models

// Модель токена, который возвращаем на клиент
// Содержит в себе access и refresh токен
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
