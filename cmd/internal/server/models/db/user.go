package db

// Пользователь, который получаем из БД
//
// Id - идентификатор пользователя
// Email - почта пользователя
//// Username - имя пользователя
// Password - пароль пользователя
type User struct {
	Id          int
	Email       string
	Password    string
}
