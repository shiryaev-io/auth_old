package routers

import (
	"auth/cmd/internal/server/controllers"
	"auth/cmd/internal/server/middlewares"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"

	"github.com/gorilla/mux"
)

const (
	get  = "GET"
	post = "POST"

	urlAuth = "/auth"
	// Путь для входа пользователя
	urlLogIn = urlAuth + "/login"
	// Путь для выхода пользователя
	urlLogOut = urlAuth + "/logout"
	// Путь для обновления токенов
	urlRefresh = urlAuth + "/refresh"
)

type ApiRouter struct {
	Router      *mux.Router
	AuthService *services.AuthService
	Logger      *logging.Logger
}

func (apiRouter *ApiRouter) Init() {

	userController := controllers.UserController{
		UserService: apiRouter.AuthService.UserService,
	}

	apiRouter.Router.HandleFunc(
		urlLogIn,
		middlewares.ErrorMiddleware(userController.Login),
	).Methods(post)
}
