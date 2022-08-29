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

// Структура, которая хранит роутер и сервисы
type ApiRouter struct {
	Router      *mux.Router
	AuthService *services.AuthService
	Logger      *logging.Logger
}

// Инициализация запросов
func (apiRouter *ApiRouter) Init() {

	userController := controllers.UserController{
		UserService: apiRouter.AuthService.UserService,
		Logger:      apiRouter.Logger,
	}
	tokenController := controllers.TokenController{
		TokenService: apiRouter.AuthService.TokenService,
		Logger:       apiRouter.Logger,
	}

	// Аутентификация пользователя
	apiRouter.handlerFunc(
		urlLogIn,
		userController.Login,
	).Methods(post)

	// Разлогин пользователя
	apiRouter.handlerFunc(
		urlLogOut,
		userController.Logout,
	).Methods(post)

	// Обновление пары access и refresh токенов
	apiRouter.handlerFunc(
		urlRefresh,
		tokenController.Refresh,
	).Methods(get)
}

// Устанавливает ErrorMiddleware
func (apiRouter *ApiRouter) handlerFunc(
	path string,
	handler middlewares.ErrorHandlerFunc,
) *mux.Route {
	return apiRouter.Router.HandleFunc(
		path,
		middlewares.ErrorMiddleware(handler).ServeHTTP,
	)
}
