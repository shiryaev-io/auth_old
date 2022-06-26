package server

import (
	"auth/pkg/logging"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	serviceHost = "SERVICE_HOST"
	servicePort = "SERVICE_PORT"
)

// Запускает сервер
func Run(router *mux.Router, logger *logging.Logger) {

	logger.Infoln(logGetEnv)
	host := os.Getenv(serviceHost)
	port := os.Getenv(servicePort)
	serviceUrl := host + ":" + port

	logger.Infof(logRunServer, serviceUrl)
	log.Fatal(http.ListenAndServe(serviceUrl, router))
}