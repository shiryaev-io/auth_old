package services

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/dtos"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models"
	"auth/cmd/pkg/logging"

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
	Logger       *logging.Logger
}

// Авторизация пользователя
func (service *UserService) Login(email, password string) (*models.Token, error) {
	service.Logger.Infof(strings.LogGettingUserByEmail, email)

	// Проверяет, существует ли пользователь в БД
	user, err := service.UserStorage.FindOne(email)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalFindUserByEmail, err)

		return nil, exceptions.BadRequest(strings.ErrorUserWithEmailNotFound, err)
	}

	service.Logger.Infoln(strings.LogCheckIfPasswordsMatch)

	// Проверяет, совпадает ли пароль,
	// который прислал клиент, с паролем, который хранится в БД
	// (пароль в БД хранится в зашифрованном виде)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	// Если пароли не совпадают
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		service.Logger.Fatalf(strings.LogFatalPasswordsNotMatch, err)

		return nil, exceptions.BadRequest(strings.ErrorInvalidPassword, err)
	}

	service.Logger.Infoln(strings.LogCreateObjectWithUserData)

	userDto := &dtos.UserDto{
		Id:          user.Id,
		Email:       user.Email,
		IsActivated: user.IsActivated,
	}

	service.Logger.Infoln(strings.LogGenerateAccessAndRefreshTokens)

	// Генерируется пара токенов: access и refresh
	tokens, err := service.TokenService.GenerateTokens(userDto)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalGenerateAccessAndRefreshTokens, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedGenerateTokens, err)
	}

	service.Logger.Infoln(strings.LogSaveRefreshTokenInDb)

	// Сохранение токена в БД
	_, err = service.TokenService.SaveToken(user.Id, tokens.RefreshToken)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalSaveRefreshTokenInDb, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedSaveRefreshToken, err)
	}

	return tokens, nil
}
