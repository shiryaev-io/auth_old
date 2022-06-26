package main

import (
	"auth/internal/server"
	"auth/internal/server/routers"
	"auth/pkg/logging"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	// "github.com/joho/godotenv"
)

func main() {
	logger := logging.GetLogger()

	// err := godotenv.Load()
	// if err != nil {
	// 	logger.Fatalf(logFailedAccessFileEnv, err)
	// } else {
	// 	logger.Fatalf(logGetEnvSuccess)
	// }

	// Канал для сигналов
	sig := make(chan bool)
	// Основной канал
	loop := make(chan error)

	// Мониторинг сигналов
	go listenerSignal(sig, logger)

	for quit := false; !quit; {
		go func() {
			initAndRanServer(logger)
			loop <- nil
		}()

		// Блокировка программы при получении сигнала
		select {
		// Прерывается выполнение программы
		case quit = <-sig:
		// Продолжается выполлнение программы
		case <-loop:
		}
	}
}

func listenerSignal(sig chan bool, logger *logging.Logger) {
	var quit bool

	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	for signal := range c {
		logger.Infof(logGetSignalSuccess, signal.String())

		switch signal {
		case syscall.SIGINT, syscall.SIGTERM:
			quit = true
		case syscall.SIGHUP:
			quit = false
		}

		if quit {
			quit = false
			// TODO: closeDB(), closeLog()
		}

		// Оповещение о прекращении работы
		sig <- quit
	}
}


func initAndRanServer(logger *logging.Logger) {

	logger.Infoln(logInitRouters)
	router := mux.NewRouter()
	routers.Init(router)

	server.Run(router, logger)
}