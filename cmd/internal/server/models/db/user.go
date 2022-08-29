package db

// Пользователь, который получаем из БД
//
// Id - идентификатор пользователя
// Email - почта пользователя
//// Username - имя пользователя
// Password - пароль пользователя
// IsActivated - подтвердил ли пользователь почту или нет
//// ActionvationLink - ссылка для активации (возможно стоит убрать)
type User struct {
	Id          int
	Email       string
	Password    string
	IsActivated bool
}
