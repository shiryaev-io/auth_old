package postgresql

import (
	"auth/cmd/internal/server/dtos"
)

// Сткруатруа для реелизации интерфейса TokenStorage
type TokenDatabase struct {
	AuthDatabase *authDatabase
}

// Получение модели токена из БД по userId
func (storage *TokenDatabase) FindByUserId(userId string) (*dtos.TokenDto, error) {
	// TODO: Реализовать метод получения токена для пользователя
	return &dtos.TokenDto{
		User:         "Test user ID",
		RefreshToken: "test.user.ID",
	}, nil
}

// Получение модели токена из БД
func (storage *TokenDatabase) FindToken(refreshToken string) (*dtos.TokenDto, error) {
	// TODO: реализовать метод получения токена
	return &dtos.TokenDto{
		User:         "Test user ID",
		RefreshToken: "test.user.ID",
	}, nil
}

// Сохранение токена в БД
func (storage *TokenDatabase) SaveToken(token *dtos.TokenDto) (*dtos.TokenDto, error) {
	// TODO: реализовать обновление токена
	return &dtos.TokenDto{
		User:         "Test user ID",
		RefreshToken: "test.user.ID",
	}, nil
}


// Создание новой записи в БД
func (storage *TokenDatabase) CreateToken(userId, refreshToken string) (*dtos.TokenDto, error) {
	// TODO: реализовать создание токена для пользователя в БД
	token := &dtos.TokenDto{
		User:         userId,
		RefreshToken: refreshToken,
	}
	return token, nil
}

// Удаляет токен из БД
func (storage *TokenDatabase) RemoveToken(refreshToken string) (*dtos.TokenDto, error) {
	// TODO: реализовать удаление токена из БД
	return &dtos.TokenDto{
		User:         "Example User",
		RefreshToken: refreshToken,
	}, nil
}
