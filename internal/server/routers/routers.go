package routers

import (
	"auth/internal/server/controllers"

	"github.com/gorilla/mux"
)

const (
	get = "GET"

	urlAuth = "/auth"
	// Путь для входа пользователя
	urlSignIn = urlAuth + "/signin"
)

// Инициализация роутеров
func Init(router *mux.Router) {
	// TODO: Для авторизации использовать POST
	router.HandleFunc(urlSignIn, controllers.SigIn).Methods(get)
}
