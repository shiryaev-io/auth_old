package strings

const (
	// Запуск сервера
	LogRunServer = "Запуск сервера на URL: %s"
	LogGetEnv    = "Получение Env"

	// Подключение базы данных
	LogAttemptConnectDb = "Попытка подключения к БД: %s"
	LogTryConnectDb     = "Подключение к БД"
	LogFatalConnectDb   = "Не удалось подключиться к БД, ошибка: %v"
	LogConnectSuccess   = "База данных успешно подключена"
	LogGetDatabaseError = "База данных не получена"

	LogGetEnvSuccess    = "Значения из файла '.env' получены."
	LogGetSignalSuccess = "Сигнал получен!"
	LogInitRouters      = "Инициализация роутера"

	// Создание токенов
	LogCreateAccessToken            = "Создание access токена"
	LogCreateRefreshToken           = "Создание refresh токена"
	LogGetJwtAccessSecret           = "Получение секретного ключа access токена из файла"
	LogGettedJwtAccessSecret        = "Секретный ключ для access токена"
	LogGetJwtRefreshSecret          = "Получение секретного ключа refresh токена из файла"
	LogGettedJwtRefreshSecret       = "Секретный ключ для refresh токена"
	LogGenerateAccessToken          = "Генерация access токена"
	LogGenerateRefreshToken         = "Генерация refresh токена"
	LogFatalGenerateAccessToken     = "Ошибка генерации access токена: %s"
	LogFatalGenerateRefreshToken    = "Ошибка генерации refresh токена: %s"
	LogSuccessGeneratedAccessToken  = "Access токен успешно был сгенерирован"
	LogSuccessGeneratedRefreshToken = "Refresh токен успешно был сгенерирован"

	// Сохранение токена
	LogGetTokenOfUser            = "Получение токена для пользователя из БД"
	LogFatalGetTokenOfUser       = "Ошибка получения токена для пользователя: %s"
	LogCreateTokenInDb           = "Создание новой записи с токеном в БД"
	LogFatalCreateTokenInDb      = "Ошибка создания новой записи с токеном в БД: %s"
	LogSuccessCreateTokenInDb    = "Новая запись с токеном успешно создана в БД"
	LogSuccesFindToken           = "Токен для пользователя %s усешно найден"
	LogUpdageRefreshToken        = "Обновление refresh токена"
	LogFatalUpdateRefreshToken   = "Ошибка обновления refresh токена: %s"
	LogSuccessUpdateRefreshToken = "Refresh токен был успешно обновлен"

	// Сценарий авторизации
	LogGettingRequestBody                  = "Получение тела запроса Body"
	LogFatalReadRequestBody                = "Не удалось прочитать тело запроса Body: %v"
	LogGettingJsonFromRequestBody          = "Получение JSON из тела запроса Body"
	LogFatalReadJsonFromRequestBody        = "Не удалось прочитать JSON из тела запроса Body: %v"
	LogUserAuthByLoginAndPassword          = "Авторизация пользователя по логину и паролю"
	LogFatalUserAuthByLoginAndPassword     = "Не удалось авторизовать пользователя по логину и паролю: %v"
	LogConvertTokensToJson                 = "Конвертация токенов в JSON"
	LogFatalConvertTokensToJson            = "Не удалось преобразовать токены в JSON: %v"
	LogGettingUserByEmail                  = "Получение пользователя по email: %s"
	LogFatalFindUserByEmail                = "Не удалось найти пользователя по email: %v"
	LogCheckIfPasswordsMatch               = "Проверка, совпадают ли пароли"
	LogFatalPasswordsNotMatch              = "Пароли не совпадают: %v"
	LogCreateObjectWithUserData            = "Создание объекта с пользовательскими данными"
	LogGenerateAccessAndRefreshTokens      = "Генерация access и refresh токенов"
	LogFatalGenerateAccessAndRefreshTokens = "Не удалось сгенерировать access и refresh токены: %v"
	LogSaveRefreshTokenInDb                = "Сохранение refresh токена в БД"
	LogFatalSaveRefreshTokenInDb           = "Не удалось сохранить refresh токена в БД: %v"

	// Сценарий разлогирования
	LogGettingRefreshTokenFromCookies = "Получение refreh токена из Cookies"
	LogFatalGettingCookies            = "Ошибка получения cookie. Необходимо передавать refresh токен: %v"
	LogFatalRefreshTokenIsEmpty       = "Refresh токен пустой"
	LogUserLogout                     = "Разлогинивание пользователя"
	LogFatalUserLogout                = "Не удалось разлогинить пользователя: %v"

	// Валидация токена
	LogStartParseAndValidateToken = "Начало парсинга и валидация токена"
	LogFatalParseJwtToken         = "Ошибка парсинга jwt токена: %v"
)
