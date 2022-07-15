package server

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	serviceHost = "SERVICE_HOST"
	servicePort = "SERVICE_PORT"
)

// Структура сервера
type server struct {
	router      *mux.Router
	authService *services.AuthService
	logger      *logging.Logger
}

// Создание сткруктуры Server
func NewServer(
	router *mux.Router,
	authService *services.AuthService,
	logger *logging.Logger,
) *server {
	return &server{
		router,
		authService,
		logger,
	}
}

// Запускает сетвер
func (server *server) Run() {
	server.logger.Infoln(strings.LogGetEnv)
	host := os.Getenv(serviceHost)
	port := os.Getenv(servicePort)
	serviceUrl := host + ":" + port

	server.logger.Infof(strings.LogRunServer, serviceUrl)
	log.Fatal(http.ListenAndServe(serviceUrl, server.router))
}
