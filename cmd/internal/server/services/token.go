package services

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/dtos"
	"auth/cmd/internal/server/models"
	"auth/cmd/pkg/logging"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtAccessSecret  = "JWT_ACCESS_SECRET"
	jwtRefreshSecret = "JWT_REFRESH_SECRET"

	accessTokenExpiresAt  = 10
	refreshTokenExpiresAt = 30
)

// Интерфейс для БД, который работает с токенами
type TokenStorage interface {
	FindOne(userId string) (*dtos.TokenDto, error)
	SaveToken(token *dtos.TokenDto) (*dtos.TokenDto, error)
	CreateToken(userId, refreshToken string) (*dtos.TokenDto, error)
}

// Сервис для токенов
type TokenService struct {
	TokenStorage TokenStorage
	Logger       *logging.Logger
}

func (service *TokenService) GenerateTokens(user *dtos.UserDto) (*models.Token, error) {
	service.Logger.Infoln(strings.LogCreateAccessToken)

	// Генерация access токена
	accessToken, err := service.createAccessToken(user)
	if err != nil {
		// TODO: вместо `error` возвращать кастомную ошибку ApiError
		return nil, err
	}

	service.Logger.Infoln(strings.LogCreateRefreshToken)

	// Генерация refresh токена
	refreshToken, err := service.createRefreshToken(user)
	if err != nil {
		// TODO: вместо `error` возвращать кастомную ошибку ApiError
		return nil, err
	}

	token := &models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return token, nil
}

// Функция сохранения refresh токена в БД
func (service *TokenService) SaveToken(userId, refreshToken string) (*dtos.TokenDto, error) {
	service.Logger.Infoln(strings.LogGetTokenOfUser)

	// Находим токен пользователя в БД
	token, err := service.TokenStorage.FindOne(userId)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalGetTokenOfUser, err)
		service.Logger.Infoln(strings.LogCreateTokenInDb)

		// Если ошибка (т.е. не нашли токен для пользователя), то создаем новую запись
		token, err = service.TokenStorage.CreateToken(
			userId,
			refreshToken,
		)
		if err != nil {
			service.Logger.Fatalf(strings.LogFatalCreateTokenInDb, err)
			return nil, err
		}

		service.Logger.Infoln(strings.LogSuccessCreateTokenInDb)

		return token, err
	}

	service.Logger.Infof(strings.LogSuccesFindToken, userId)
	service.Logger.Infoln(strings.LogUpdageRefreshToken)

	// Обновляем refresh токен у пользователя
	token.RefreshToken = refreshToken
	token, err = service.TokenStorage.SaveToken(token)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalUpdateRefreshToken, err)
		return nil, err
	}

	service.Logger.Infoln(strings.LogSuccessUpdateRefreshToken)

	return token, err
}

// Создание access токена
func (service *TokenService) createAccessToken(user *dtos.UserDto) (string, error) {
	// время, черз которое access токен протухнет
	expiredAt := time.Now().Add(
		time.Minute * accessTokenExpiresAt,
	).Unix()

	claims := jwt.StandardClaims{
		IssuedAt:  expiredAt,
		Subject:   user.Id,
	}

	service.Logger.Infoln(strings.LogGetJwtAccessSecret)
	signJwtAceessSecret := os.Getenv(jwtAccessSecret)
	service.Logger.Infof(strings.LogGettedJwtAccessSecret, signJwtAceessSecret)

	service.Logger.Infoln(strings.LogGenerateAccessToken)

	// Генерация access токена
	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(
		[]byte(signJwtAceessSecret), // Токен подписывается секретным ключом из .env
	)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalGenerateAccessToken, err)
		return strings.Empty, err
	}

	service.Logger.Infoln(strings.LogSuccessGeneratedAccessToken)

	return accessToken, nil
}

// Создание refresh токена
func (service *TokenService) createRefreshToken(user *dtos.UserDto) (string, error) {
	// время, черз которое refresh токен протухнет
	expiredAt := time.Now().Add(
		time.Minute * refreshTokenExpiresAt,
	).Unix()

	claims := jwt.StandardClaims{
		IssuedAt:  expiredAt,
		Subject:   user.Id,
	}

	service.Logger.Infoln(strings.LogGetJwtRefreshSecret)
	signJwtRefreshSecret := os.Getenv(jwtRefreshSecret)
	service.Logger.Infoln(strings.LogGettedJwtRefreshSecret, signJwtRefreshSecret)

	service.Logger.Infoln(strings.LogGenerateRefreshToken)

	// Генерация refresh токена
	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(
		[]byte(signJwtRefreshSecret), // Токен подписывается секретным ключом из .env
	)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalGenerateRefreshToken, err)
		return strings.Empty, err
	}

	service.Logger.Infoln(strings.LogSuccessGeneratedRefreshToken)

	return refreshToken, nil
}
