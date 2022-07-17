package dtos

// Dto модель пользователя БЕЗ пароля
// (т.к. будем отправлять на клиент)
type UserDto struct {
	Id          string
	Email       string
	IsActivated bool
}
