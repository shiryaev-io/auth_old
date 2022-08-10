package services

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/dtos"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models"
	"auth/cmd/pkg/logging"
	"errors"
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
	FindByUserId(userId string) (*dtos.TokenDto, error)
	FindToken(refreshToken string) (*dtos.TokenDto, error)
	SaveToken(token *dtos.TokenDto) (*dtos.TokenDto, error)
	CreateToken(userId, refreshToken string) (*dtos.TokenDto, error)
	RemoveToken(refreshToken string) (*dtos.TokenDto, error)
}

// Сервис для токенов
type TokenService struct {
	TokenStorage TokenStorage
	UserStorage  UserStorage
	Logger       *logging.Logger
}

func (service *TokenService) GenerateTokens(user *dtos.UserDto) (*models.Token, error) {
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
	token, err := service.TokenStorage.FindByUserId(userId)
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

// Обновляет пару access и refresh токенов
func (service *TokenService) Refresh(refreshToken string) (*models.Token, error) {
	service.Logger.Infoln(strings.LogGetJwtRefreshSecret)

	signJwtRefreshSecret := os.Getenv(jwtRefreshSecret)

	service.Logger.Infoln("Парсинг и валидация токена")

	payload, errorValidate := service.validateToken(refreshToken, signJwtRefreshSecret)

	service.Logger.Infoln("Каст интерфейса Claims в структуру StandardClaims")

	standardClaims := payload.(jwt.StandardClaims)

	service.Logger.Infoln("Поиск refresh токена в БД")

	_, errorFindToken := service.TokenStorage.FindToken(refreshToken)
	if errorValidate != nil || errorFindToken != nil {
		service.Logger.Fatalf("Не удалось найти или провалидировать токен: Ошибка 1: %v; Ошибка 2: %v", errorValidate, errorFindToken)

		return nil, exceptions.UnauthorizedError(errorValidate)
	}

	service.Logger.Infoln("Поиск пользователя по ID в БД")

	user, err := service.UserStorage.FindById(standardClaims.Id)
	if err != nil {
		service.Logger.Fatalf("Не удалось найти пользователя по ID в БД", err)

		return nil, exceptions.BadRequest(strings.ErrorUserWithEmailNotFound, err)
	}

	userDto := &dtos.UserDto{
		Id:          user.Id,
		Email:       user.Email,
		IsActivated: user.IsActivated,
	}

	service.Logger.Infoln(strings.LogGenerateAccessAndRefreshTokens)

	// Генерируется пара токенов: access и refresh
	tokens, err := service.GenerateTokens(userDto)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalGenerateAccessAndRefreshTokens, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedGenerateTokens, err)
	}

	service.Logger.Infoln(strings.LogSaveRefreshTokenInDb)

	// Сохранение токена в БД
	_, err = service.SaveToken(user.Id, tokens.RefreshToken)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalSaveRefreshTokenInDb, err)

		return nil, exceptions.BadRequest(strings.ErrorFailedSaveRefreshToken, err)
	}

	return tokens, nil
}

// Создание access токена
func (service *TokenService) createAccessToken(user *dtos.UserDto) (string, error) {
	// время, черз которое access токен протухнет
	expiredAt := time.Now().Add(
		time.Minute * accessTokenExpiresAt,
	).Unix()

	claims := jwt.StandardClaims{
		ExpiresAt: expiredAt,
		Subject:   user.Id,
	}

	service.Logger.Infoln(strings.LogGetJwtAccessSecret)
	signJwtAceessSecret := os.Getenv(jwtAccessSecret)
	service.Logger.Infof(strings.LogGettedJwtAccessSecret, signJwtAceessSecret)

	service.Logger.Infoln(strings.LogGenerateAccessToken)

	// Генерация access токена
	accessToken, err := service.createJwt(claims, signJwtAceessSecret)
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
		ExpiresAt: expiredAt,
		Subject:   user.Id,
	}

	service.Logger.Infoln(strings.LogGetJwtRefreshSecret)
	signJwtRefreshSecret := os.Getenv(jwtRefreshSecret)
	service.Logger.Infoln(strings.LogGettedJwtRefreshSecret, signJwtRefreshSecret)

	service.Logger.Infoln(strings.LogGenerateRefreshToken)

	// Генерация refresh токена
	refreshToken, err := service.createJwt(claims, signJwtRefreshSecret)
	if err != nil {
		service.Logger.Fatalf(strings.LogFatalGenerateRefreshToken, err)
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

// Валидация токена. Необходимо, чтобы понимать, что токен не был подделан
// или что срок годности не иссяк
func (service *TokenService) validateToken(tokenString, signingKey string) (jwt.Claims, error) {
	service.Logger.Infoln("Начало парсинг и валидация токена")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// TODO: вынести строку в ресурсы
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		service.Logger.Fatalf("Ошибка парсинга jwt токена: %v", err)

		return nil, err
	}
	return token.Claims, nil
}

// Удаление одного токена из БД
func (service *TokenService) RemoveToken(refreshToken string) (*dtos.TokenDto, error) {
	return service.TokenStorage.RemoveToken(refreshToken)
}
