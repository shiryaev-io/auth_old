package middlewares

import (
	stringsRes "auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const (
	authorizationHeader = "Authorization"
	bearer              = "Bearer"
	userDataKey         = "userData"
)

// Содержит сервис для работы с токенами
type AuthMiddleware struct {
	TokenService *services.TokenService
	Logger       *logging.Logger
}

// Авторизовывает пользователя по access токену
func (middleware *AuthMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		payload, err := middleware.authorizationByToken(request)
		if err != nil {
			handleError(response, exceptions.UnauthorizedError(err))
			return
		}

		// TODO: в контекст передавать данные из payload
		standardClaims := payload.(*jwt.StandardClaims)
		// TODO: обернуть userId в стркутуру (типа UserRequest)
		userId := standardClaims.Id
		parentContext := request.Context()
		contextWithUserData := context.WithValue(parentContext, userDataKey, userId)
		requestWithContext := request.Clone(contextWithUserData)
		handler.ServeHTTP(response, requestWithContext)
	})
}

// Авторизовывает пользователя по header Authorization
func (middleware *AuthMiddleware) authorizationByToken(request *http.Request) (jwt.Claims, error) {
	header := request.Header
	authHeader := header.Get(authorizationHeader)
	if authHeader == stringsRes.Empty {
		return nil, errors.New(stringsRes.ErrorHeaderAuthorizationIsEmpty)
	}

	headerParts := strings.Split(authHeader, stringsRes.Space)
	if len(headerParts) != 2 || headerParts[0] != bearer {
		return nil, errors.New(stringsRes.ErrorInvalidAuthorizationHeader)
	}

	accessToken := headerParts[1]
	if len(accessToken) == 0 {
		return nil, errors.New(stringsRes.ErrorTokenIsEmpty)
	}

	payload, err := middleware.TokenService.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
