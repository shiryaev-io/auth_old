package strings

const (
	Empty = ""

	// Строки ошибок
	ErrorUserWithEmailNotFound      = "Пользователь с таким email не найден"
	ErrorInvalidPassword            = "Неверный пароль"
	ErrorFailedGenerateTokens       = "Не удалось сгенерировать токены"
	ErrorFailedSaveRefreshToken     = "Не удалось сохранить refresh токен"
	ErrorFailedLogin                = "Не удалось авторизоваться, попробуйте позже"
	ErrorInvalidData                = "Неверные данные, попробуйте еще раз"
	ErrorWrongLoginOrPassword       = "Неверный логин или пароль. Попробуйте еще раз"
	ErrorInternal                   = "Внутренняя ошибка. Попробуйте позже"
	ErrorTryAgaint                  = "Ошибка, попробуйте еще раз"
	ErrorRefreshTokenMustNotBeEmpty = "Refresh токен не должен быть пустым"
	ErrorLogout                     = "Ошибка разлогина. Попробуйте еще раз"
)
