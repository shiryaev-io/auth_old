package controllers

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models/requests"
	"auth/cmd/internal/server/models/responses"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	cookieRefreshToken = "refreshToken"
)

// Контроллер пользователя
// Необходим для работы с пользователями, 
// например авторизация, разлогирование и т.д.
type UserController struct {
	UserService *services.UserService
	Logger      *logging.Logger
}

// Авторизация пользователя
func (controller *UserController) Login(
	response http.ResponseWriter,
	request *http.Request,
) (*responses.Common, error) {
	controller.Logger.Infoln(strings.LogGettingRequestBody)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalReadRequestBody, err)

		return nil, exceptions.BadRequest(strings.ErrorInvalidData, err)
	}

	controller.Logger.Infoln(strings.LogGettingJsonFromRequestBody)

	userRequest := requests.UserRequest{}
	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalReadJsonFromRequestBody, err)
		
		return nil, exceptions.BadRequest(strings.ErrorInvalidData, err)
	}

	controller.Logger.Infoln(strings.LogUserAuthByLoginAndPassword)

	tokens, err := controller.UserService.Login(
		userRequest.Email,
		userRequest.Password,
	)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalUserAuthByLoginAndPassword, err)
		
		return nil, exceptions.BadRequest(strings.ErrorWrongLoginOrPassword, err)
	}

	controller.Logger.Infoln(strings.LogConvertTokensToJson)

	jsonBody, err := json.Marshal(tokens)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalConvertTokensToJson, err)
		
		return nil, exceptions.ServerError(strings.ErrorInternal, err)
	}
	return &responses.Common{
		Status: http.StatusOK,
		Body: jsonBody,
	}, nil
}

// Разлогин пользователя
func (controller *UserController) Logout(
	response http.ResponseWriter,
	request *http.Request,
) (*responses.Common, error) {
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

	controller.Logger.Infoln(strings.LogUserLogout)

	err = controller.UserService.Logout(refreshToken)
	if err != nil {
		controller.Logger.Infof(strings.LogFatalUserLogout, err)

		return nil, exceptions.BadRequest(strings.ErrorLogout, err)
	}

	return &responses.Common{
		Status: http.StatusNoContent,
		Body:   nil,
	}, nil
}
