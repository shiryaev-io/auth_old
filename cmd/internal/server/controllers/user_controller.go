package controllers

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/requests"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	cookieRefreshToken = "refresh_token"
)

// Контроллер пользователя
type UserController struct {
	UserService *services.UserService
	Logger      *logging.Logger
}

// Авторизация пользователя
func (controller *UserController) Login(
	response http.ResponseWriter,
	request *http.Request,
) error {
	controller.Logger.Infoln(strings.LogGettingRequestBody)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		controller.Logger.Fatalf(strings.LogFatalReadRequestBody, err)

		return exceptions.BadRequest(strings.ErrorInvalidData, err)
	}

	controller.Logger.Infoln(strings.LogGettingJsonFromRequestBody)

	userRequest := requests.UserRequest{}
	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		controller.Logger.Fatalf(strings.LogFatalReadJsonFromRequestBody, err)
		
		return exceptions.BadRequest(strings.ErrorInvalidData, err)
	}

	controller.Logger.Infoln(strings.LogUserAuthByLoginAndPassword)

	tokens, err := controller.UserService.Login(
		userRequest.Email,
		userRequest.Password,
	)
	if err != nil {
		controller.Logger.Fatalf(strings.LogFatalUserAuthByLoginAndPassword, err)
		
		return exceptions.BadRequest(strings.ErrorWrongLoginOrPassword, err)
	}

	controller.Logger.Infoln(strings.LogConvertTokensToJson)

	jsonBody, err := json.Marshal(tokens)
	if err != nil {
		controller.Logger.Fatalf(strings.LogFatalConvertTokensToJson, err)
		
		return exceptions.ServerError(strings.ErrorInternal, err)
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonBody)
	return nil
}

// Разлогин пользователя
func (controller *UserController) Logout(
	response http.ResponseWriter,
	request *http.Request,
) error {
	controller.Logger.Infoln(strings.LogGettingRefreshTokenFromCookies)

	cookie, err := request.Cookie(cookieRefreshToken)
	if err != nil {
		controller.Logger.Fatalf(strings.LogFatalGettingCookies, err)

		return exceptions.BadRequest(strings.ErrorTryAgaint, err)
	}

	refreshToken := cookie.Value
	if refreshToken == strings.Empty {
		controller.Logger.Fatalf(strings.LogFatalRefreshTokenIsEmpty, err)

		return exceptions.BadRequest(strings.ErrorRefreshTokenMustNotBeEmpty, err)
	}

	controller.Logger.Infoln(strings.LogUserLogout)

	_, err = controller.UserService.Logout(refreshToken)
	if err != nil {
		controller.Logger.Fatalf(strings.LogFatalUserLogout, err)

		return exceptions.BadRequest(strings.ErrorLogout, err)
	}

	response.WriteHeader(http.StatusNoContent)

	return nil
}
