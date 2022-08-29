package controllers

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models/responses"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"
	"encoding/json"
	"net/http"
)

// Контроллер для токенов
// Необходим для работы с токенами,
// например для обновления пары access и refresh токенов
type TokenController struct {
	TokenService *services.TokenService
	Logger       *logging.Logger
}

// Обновляет пару access и refresh токенов
func (controller *TokenController) Refresh(
	response http.ResponseWriter,
	request *http.Request,
) (*responses.Common, error) {
	// TODO: подумать, возможно получение refresh токена из cookie вынести в общий код,
	// т.к сейчас мы достаем refresh токен еще и в UserController в  Logout
	// start region: Получение refresh токена
	controller.Logger.Infoln(strings.LogGettingRefreshTokenFromCookies)

	cookie, err := request.Cookie(cookieRefreshToken)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalGettingCookies, err)

		return nil, exceptions.BadRequest(strings.ErrorTryAgaint, err)
	}

	refreshToken := cookie.Value
	if refreshToken == strings.Empty {
		controller.Logger.Infof(strings.LogFatalRefreshTokenIsEmpty, err)

		return nil, exceptions.UnauthorizedError(err)
	}
	// end region: Получение refresh токена

	tokens, err := controller.TokenService.Refresh(refreshToken)
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(tokens)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalConvertTokensToJson, err)

		return nil, exceptions.ServerError(strings.ErrorInternal, err)
	}

	return &responses.Common{
		Status: http.StatusOK,
		Body:   jsonBody,
	}, nil
}
