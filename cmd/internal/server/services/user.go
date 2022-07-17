package services

import (
	"auth/cmd/internal/server/dtos"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models"

	"golang.org/x/crypto/bcrypt"
)

// Хранилище для пользователей
type UserStorage interface {
	FindOne(email string) (*models.User, error)
}

// Сервис для работы с пользователями
type UserService struct {
	UserStorage  UserStorage
	TokenService *TokenService
}

// Авторизация пользователя
func (service *UserService) Login(email, password string) (*models.Token, error) {

	// Проверяет, существует ли пользователь в БД
	user, err := service.UserStorage.FindOne(email)
	if err != nil {
		// TODO: Вынести строку в strings
		return nil, exceptions.BadRequest("Пользователь с таким email не найден", err)
	}
	// Проверяет, совпадает ли пароль,
	// который прислал клиент, с паролем, который хранится в БД
	// (пароль в БД хранится в зашифрованном виде)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	// Если пароли не совпадают
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// TODO: Вынести строку в strings
		return nil, exceptions.BadRequest("Неверный пароль", err)
	}

	userDto := &dtos.UserDto{
		Id:          user.Id,
		Email:       user.Email,
		IsActivated: user.IsActivated,
	}
	// Генерируется пара токенов: access и refresh
	tokens, err := service.TokenService.GenerateTokens(userDto)
	if err != nil {
		// TODO: Вынести строку в strings
		return nil, exceptions.BadRequest("Не удалось сгенерировать токен", err)
	}

	// Сохранение токена в БД
	_, err = service.TokenService.SaveToken(user.Id, tokens.RefreshToken)
	if err != nil {
		// TODO: Вынести строку в strings
		return nil, exceptions.BadRequest("Не удалось сохранить токен", err)
	}
	return tokens, nil
}
