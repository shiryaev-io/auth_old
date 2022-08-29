package db

// Модель токена, которую получаем из БД
//
// Id - Идентификатор токена в БД
// UserId - Идентификатор пользователя, с которым связан токен
// Value - Сам токен
type Token struct {
	Id     int
	UserId int
	Value  string
}
