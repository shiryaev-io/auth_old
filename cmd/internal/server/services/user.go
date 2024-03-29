package services

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models/db"
	"auth/cmd/internal/server/models/dto"
	"auth/cmd/pkg/logging"

	"golang.org/x/crypto/bcrypt"
)

// Хранилище для пользователей
type UserStorage interface {
	FindByEmail(email string) (*db.User, error)
	FindById(userId int) (*db.User, error)
}

// Сервис для работы с пользователями
type UserService struct {
	UserStorage  UserStorage
	TokenService *TokenService
	Logger       *logging.Logger
}

// Авторизация пользователя
func (service *UserService) Login(email, password string) (*dto.Tokens, error) {
	service.Logger.Infof(strings.LogGettingUserByEmail, email)

	// Проверяет, существует ли пользователь в БД
	user, err := service.UserStorage.FindByEmail(email)
	if err != nil {
		service.Logger.Infof(strings.LogFatalFindUserByEmail, err)

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
		service.Logger.Infof(strings.LogFatalPasswordsNotMatch, err)

		return nil, exceptions.BadRequest(strings.ErrorInvalidPassword, err)
	}

	service.Logger.Infoln(strings.LogCreateObjectWithUserData)

	userDto := &dto.User{
		Id:          user.Id,
		Email:       user.Email,
	}

	service.Logger.Infoln(strings.LogGenerateAccessAndRefreshTokens)

	// Генерируется пара токенов: access и refresh
	tokens, err := service.TokenService.GenerateTokens(userDto)
	if err != nil {
		service.Logger.Infof(strings.LogFatalGenerateAccessAndRefreshTokens, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedGenerateTokens, err)
	}

	service.Logger.Infoln(strings.LogSaveRefreshTokenInDb)

	// Сохранение токена в БД
	_, err = service.TokenService.SaveToken(user.Id, tokens.Refresh)
	if err != nil {
		service.Logger.Infof(strings.LogFatalSaveRefreshTokenInDb, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedSaveRefreshToken, err)
	}

	return tokens, nil
}

// Разлогин пользователя
func (service *UserService) Logout(refreshToken string) error {
	service.Logger.Infoln(strings.LogCallingRefreshTokenDeletaionFun)

	err := service.removeToken(refreshToken)
	if err != nil {
		service.Logger.Infof(strings.LogFatalDeleteRefreshToken, err)

		return err
	}

	service.Logger.Infoln(strings.LogRefreshTokenSuccessDeleted)

	return nil
}

// Логика удаления токена из БД
func (service *UserService) removeToken(refreshToken string) error {
	service.Logger.Infoln(strings.LogCallingTokenRemovalFunFromDb)

	err := service.TokenService.RemoveToken(refreshToken)
	if err != nil {
		service.Logger.Infof(strings.LogFatalGetRefreshTokenFromDb, err)

		return err
	}

	service.Logger.Infoln(strings.LogRefreshTokenSuccessDeletedFromDb)

	return nil
}
