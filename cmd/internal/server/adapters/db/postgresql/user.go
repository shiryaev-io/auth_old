package postgresql

import "auth/cmd/internal/server/models"

// Структура для реализации интрефейса UserStorage
type UserDatabase struct {
	AuthDatabase *authDatabase
}

func (storage *UserDatabase) FindOne(email string) (*models.User, error) {
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