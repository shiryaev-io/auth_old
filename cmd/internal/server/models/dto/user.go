package dto

// Данные пользователя, которые передаются между слоями
//
// Id - идентификатор пользователя
// Email - почта пользователя
//// Username - имя пользователя
type User struct {
	Id          int
	Email       string
}
