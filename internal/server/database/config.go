package database

import "os"

// Даные для подключения к БД
type dbConfig struct {
	driver   string
	host     string
	port     string
	name     string
	user     string
	password string
}

// Получение конфигурации БД
func getDbConfig() *dbConfig {
	return &dbConfig{
		driver:   os.Getenv(dbDriver),
		host:     os.Getenv(dbHost),
		port:     os.Getenv(dbPort),
		name:     os.Getenv(dbName),
		user:     os.Getenv(dbUser),
		password: os.Getenv(dbPassword),
	}
}
