package dto

// Данные пользователя, которые передаются между слоями
//
// Id - идентификатор пользователя
// Email - почта пользователя
//// Username - имя пользователя
// IsActivated - подтвердил ли пользователь почту или нет
//// ActionvationLink - ссылка для активации (возможно стоит убрать)
type User struct {
	Id          int
	Email       string
	IsActivated bool
}
