package routers

import (
	"auth/internal/server/controllers"

	"github.com/gorilla/mux"
)

const (
	urlAuth = "/auth"

	// Путь для входа пользователя
	urlSignIn = urlAuth + "/signin"
)

// Инициализация роутеров
func Init(router *mux.Router) {
	router.HandleFunc(urlSignIn, controllers.SigIn).Methods("POST")
}
