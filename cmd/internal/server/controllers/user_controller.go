package controllers

import (
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/requests"
	"auth/cmd/internal/server/services"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Авторизация пользователя
// TODO: удалить этот метод
func LogIn(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hello, world"))
}

// Контроллер пользователя
type UserController struct {
	UserService *services.UserService
}

// Контроллер для входа пользователя
func (controller *UserController) Login(
	response http.ResponseWriter,
	request *http.Request,
) error {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		// TODO: Вынести строку в strings
		return exceptions.BadRequest("Неверные данные", err)
	}
	userRequest := requests.UserRequest{}
	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		// TODO: Вынести строку в strings
		return exceptions.BadRequest("Неверные данные", err)
	}

	tokens, err := controller.UserService.Login(
		userRequest.Email,
		userRequest.Password,
	)
	if err != nil {
		// TODO: Вынести строку в strings
		return exceptions.BadRequest("Неверные данные", err)
	}

	jsonBody, err := json.Marshal(tokens)
	if err != nil {
		// TODO: Вынести строку в strings
		return exceptions.BadRequest("Внутренняя ошибка", err)
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonBody)
	return nil
}
