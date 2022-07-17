package models

// Модель пользователя
// Id - идентификатор пользователя
// Email - почта пользователя
// Username - имя пользователя
// Password - пароль пользователя
// IsActivated - подтвердил ли пользователь почту или нет
// ActionvationLink - ссылка для активации (возможно стоит убрать)
type User struct {
	Id              string
	Email           string
	Username        string
	Password        string
	IsActivated     bool
	ActiovationLink string
}
