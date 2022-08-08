package strings

const (
	Empty = ""

	// Строки ошибок
	ErrorUserWithEmailNotFound  = "Пользователь с таким email не найден"
	ErrorInvalidPassword        = "Неверный пароль"
	ErrorFailedGenerateTokens   = "Не удалось сгенерировать токены"
	ErrorFailedSaveRefreshToken = "Не удалось сохранить refresh токен"
	ErrorFailedLogin            = "Не удалось авторизоваться, попробуйте позже"
)
