package config

import "os"

const (
	dbDriver   = "DB_DRIVER"
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbName     = "POSTGRES_DB"
	dbUser     = "POSTGRES_USER"
	dbPassword = "POSTGRES_PASSWORD"
)

// Конфигурация БД
type ConfigDb struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func NewConfigDb() *ConfigDb {
	return &ConfigDb{
		Driver:   os.Getenv(dbDriver),
		Host:     os.Getenv(dbHost),
		Port:     os.Getenv(dbPort),
		Name:     os.Getenv(dbName),
		User:     os.Getenv(dbUser),
		Password: os.Getenv(dbPassword),
	}
}
