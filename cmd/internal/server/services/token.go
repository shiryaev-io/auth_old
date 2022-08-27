package services

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models/db"
	"auth/cmd/internal/server/models/dto"
	"auth/cmd/pkg/logging"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtAccessSecret  = "JWT_ACCESS_SECRET"
	jwtRefreshSecret = "JWT_REFRESH_SECRET"

	// TODO: получать время из файла .env
	accessTokenExpiresAt  = 5
	refreshTokenExpiresAt = 15
)

// Интерфейс для БД, который работает с токенами
type TokenStorage interface {
	FindByUserId(userId int) (*db.Token, error)
	FindToken(refreshToken string) (*db.Token, error)
	SaveToken(token *dto.Token) error
	CreateToken(userId int, refreshToken string) (*db.Token, error)
	RemoveToken(refreshToken string) (*db.Token, error)
}

// Сервис для токенов
type TokenService struct {
	TokenStorage TokenStorage
	UserStorage  UserStorage
	Logger       *logging.Logger
}

func (service *TokenService) GenerateTokens(user *dto.User) (*dto.Tokens, error) {
	service.Logger.Infoln(strings.LogCreateAccessToken)

	// Генерация access токена
	accessToken, err := service.createAccessToken(user)
	if err != nil {
		return nil, exceptions.ServerError(strings.ErrorFailedLogin, err)
	}

	service.Logger.Infoln(strings.LogCreateRefreshToken)

	// Генерация refresh токена
	refreshToken, err := service.createRefreshToken(user)
	if err != nil {
		return nil, exceptions.ServerError(strings.ErrorFailedLogin, err)
	}

	tokens := &dto.Tokens{
		Access: accessToken,
		Refresh: refreshToken,
	}

	return tokens, nil
}

// Функция сохранения refresh токена в БД
func (service *TokenService) SaveToken(userId int, refreshToken string) (*dto.Token, error) {
	service.Logger.Infoln(strings.LogGetTokenOfUser)

	// Находим токен пользователя в БД
	tokenFromDb, err := service.TokenStorage.FindByUserId(userId)
	if err != nil {
		service.Logger.Infof(strings.LogFatalGetTokenOfUser, err)
		service.Logger.Infoln(strings.LogCreateTokenInDb)

		// Если ошибка (т.е. не нашли токен для пользователя), то создаем новую запись
		tokenFromDb, err = service.TokenStorage.CreateToken(
			userId,
			refreshToken,
		)
		if err != nil {
			service.Logger.Infof(strings.LogFatalCreateTokenInDb, err)
			return nil, err
		}

		service.Logger.Infoln(strings.LogSuccessCreateTokenInDb)

		tokenDto := &dto.Token{
			Id:     tokenFromDb.Id,
			UserId: tokenFromDb.UserId,
			Value:  tokenFromDb.Value,
		}

		return tokenDto, err
	}

	service.Logger.Infof(strings.LogSuccesFindToken, userId)
	service.Logger.Infoln(strings.LogUpdageRefreshToken)

	// Обновляем refresh токен у пользователя
	tokenDto := &dto.Token{
		Id:     tokenFromDb.Id,
		UserId: tokenFromDb.UserId,
		Value:  refreshToken,
	}
	err = service.TokenStorage.SaveToken(tokenDto)
	if err != nil {
		service.Logger.Infof(strings.LogFatalUpdateRefreshToken, err)
		return nil, err
	}

	service.Logger.Infoln(strings.LogSuccessUpdateRefreshToken)

	return tokenDto, err
}

// Обновляет пару access и refresh токенов
func (service *TokenService) Refresh(refreshToken string) (*dto.Tokens, error) {
	// TODO: вынести в строки
	service.Logger.Infoln("Валидация Refresh токена")

	payload, err := service.validateRefreshToken(refreshToken)
	if err != nil {
		service.Logger.Infof("Ошибка валидации токена: %v", err)

		return nil, exceptions.UnauthorizedError(err)
	}

	service.Logger.Infoln("Каст интерфейса Claims в структуру StandardClaims")

	standardClaims := payload.(*jwt.StandardClaims)

	service.Logger.Infoln("Поиск refresh токена в БД")

	_, err = service.TokenStorage.FindToken(refreshToken)
	if err != nil {
		service.Logger.Infof("Не удалось найти токен в БД: %v", err)

		return nil, exceptions.UnauthorizedError(err)
	}

	userId, err := strconv.Atoi(standardClaims.Subject)
	if err != nil {
		return nil, exceptions.BadRequest("Тип id пользователя должен быть Int", err)
	}

	user, err := service.UserStorage.FindById(userId)
	if err != nil {
		service.Logger.Infof("Не удалось найти пользователя по ID в БД", err)

		return nil, exceptions.BadRequest(strings.ErrorUserWithEmailNotFound, err)
	}

	userDto := &dto.User{
		Id:          user.Id,
		Email:       user.Email,
		IsActivated: user.IsActivated,
	}

	service.Logger.Infoln(strings.LogGenerateAccessAndRefreshTokens)

	// Генерируется пара токенов: access и refresh
	tokens, err := service.GenerateTokens(userDto)
	if err != nil {
		service.Logger.Infof(strings.LogFatalGenerateAccessAndRefreshTokens, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedGenerateTokens, err)
	}

	service.Logger.Infoln(strings.LogSaveRefreshTokenInDb)

	// Сохранение токена в БД
	_, err = service.SaveToken(user.Id, tokens.Refresh)
	if err != nil {
		service.Logger.Infof(strings.LogFatalSaveRefreshTokenInDb, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedSaveRefreshToken, err)
	}

	return tokens, nil
}

// Создание access токена
func (service *TokenService) createAccessToken(user *dto.User) (string, error) {
	// время, черз которое access токен протухнет
	expiredAt := time.Now().Add(
		time.Minute * accessTokenExpiresAt,
	).Unix()

	claims := jwt.StandardClaims{
		ExpiresAt: expiredAt,
		Subject:   fmt.Sprint(user.Id),
	}

	service.Logger.Infoln(strings.LogGetJwtAccessSecret)
	signJwtAceessSecret := os.Getenv(jwtAccessSecret)
	service.Logger.Infoln(strings.LogGettedJwtAccessSecret)

	service.Logger.Infoln(strings.LogGenerateAccessToken)

	// Генерация access токена
	accessToken, err := service.createJwt(claims, signJwtAceessSecret)
	if err != nil {
		service.Logger.Infof(strings.LogFatalGenerateAccessToken, err)
		return strings.Empty, err
	}

	service.Logger.Infoln(strings.LogSuccessGeneratedAccessToken)

	return accessToken, nil
}

// Создание refresh токена
func (service *TokenService) createRefreshToken(user *dto.User) (string, error) {
	// время, черз которое refresh токен протухнет
	expiredAt := time.Now().Add(
		time.Minute * refreshTokenExpiresAt,
	).Unix()

	claims := jwt.StandardClaims{
		ExpiresAt: expiredAt,
		Subject:   fmt.Sprint(user.Id),
	}

	service.Logger.Infoln(strings.LogGetJwtRefreshSecret)
	signJwtRefreshSecret := os.Getenv(jwtRefreshSecret)
	service.Logger.Infoln(strings.LogGettedJwtRefreshSecret)

	service.Logger.Infoln(strings.LogGenerateRefreshToken)

	// Генерация refresh токена
	refreshToken, err := service.createJwt(claims, signJwtRefreshSecret)
	if err != nil {
		service.Logger.Infof(strings.LogFatalGenerateRefreshToken, err)
		return strings.Empty, err
	}

	service.Logger.Infoln(strings.LogSuccessGeneratedRefreshToken)

	return refreshToken, nil
}

// Генерация JWT токена
func (service *TokenService) createJwt(
	claims jwt.StandardClaims,
	secretKey string,
) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(
		[]byte(secretKey), // Токен подписывается секретным ключом
	)
}

// Валидация access токена
func (service *TokenService) ValidateAccessToken(tokenString string) (jwt.Claims, error) {
	signJwtAceessSecret := os.Getenv(jwtAccessSecret)
	payload, err := service.validateToken(tokenString, signJwtAceessSecret)
	return payload, err
}

// Валидация access токена
func (service *TokenService) validateRefreshToken(tokenString string) (jwt.Claims, error) {
	signJwtRefreshSecret := os.Getenv(jwtRefreshSecret)
	payload, err := service.validateToken(tokenString, signJwtRefreshSecret)
	return payload, err
}

// Валидация токена. Необходимо, чтобы понимать, что токен не был подделан
// или что срок годности не иссяк
func (service *TokenService) validateToken(tokenString, signingKey string) (jwt.Claims, error) {
	service.Logger.Infoln(strings.LogStartParseAndValidateToken)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(strings.ErrorUnexpectedSigningMethod)
			}
			return []byte(signingKey), nil
		},
	)
	if err != nil {
		// TODO: вынести строку в ресурсы
		service.Logger.Infof(strings.LogFatalParseJwtToken, err)

		return nil, err
	}
	return token.Claims, nil
}

// Удаление одного токена из БД
func (service *TokenService) RemoveToken(refreshToken string) (*dto.Token, error) {
	tokenFromDb, err := service.TokenStorage.RemoveToken(refreshToken)
	tokenDto := &dto.Token{
		Id:     tokenFromDb.Id,
		UserId: tokenFromDb.UserId,
		Value:  tokenFromDb.Value,
	}
	return tokenDto, err
}
