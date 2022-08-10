package postgresql

import (
	"auth/cmd/internal/server/models"
)

// Структура для реализации интрефейса UserStorage
type UserDatabase struct {
	AuthDatabase *authDatabase
}

// Получение данных пользователя по Email
func (storage *UserDatabase) FindByEmail(email string) (*models.User, error) {
	// TODO: Релизовать метод получения пользователя по email 
	return &models.User{
		Id:              "Test user ID",
		Email:           "test.user@email.test",
		Username:        "Test username",
		Password:        "Test password",
		IsActivated:     true,
		ActiovationLink: "test link",
	}, nil
}

// Получение данных пользователя по ID
func (storage *UserDatabase) FindById(userId string) (*models.User, error) {
	// TODO: Релизовать метод получения пользователя по ID 
	return &models.User{
		Id:              "Test user ID",
		Email:           "test.user@email.test",
		Username:        "Test username",
		Password:        "Test password",
		IsActivated:     true,
		ActiovationLink: "test link",
	}, nil
}
