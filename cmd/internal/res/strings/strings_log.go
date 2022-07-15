package strings

const (
	// Логи запуска сервера
	LogRunServer = "Запуск сервера на URL: %s"
	LogGetEnv    = "Получение Env"

	// Логи подключения базы данных
	LogAttemptConnectDb = "Попытка подключения к БД: %s"
	LogTryConnectDb     = "Подключение к БД"
	LogFatalConnectDb   = "Не удалось подключиться к БД, ошибка: %v"
	LogConnectSuccess   = "База данных успешно подключена"
	LogGetDatabaseError = "База данных не получена"

	LogGetEnvSuccess    = "Значения из файла '.env' получены."
	LogGetSignalSuccess = "Сигнал получен!"
	LogInitRouters      = "Инициализация роутера"
)
