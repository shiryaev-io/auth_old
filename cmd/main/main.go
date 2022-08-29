package main

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server"
	"auth/cmd/internal/server/adapters/db/postgresql"
	"auth/cmd/internal/server/config"
	"auth/cmd/internal/server/routers"
	"auth/cmd/internal/server/services"
	"auth/cmd/pkg/logging"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	logger := logging.GetLogger()

	// Канал для сигналов
	sig := make(chan bool)
	// Основной канал
	loop := make(chan error)

	// Мониторинг сигналов
	go listenerSignal(sig, logger)

	for quit := false; !quit; {
		go func() {
			initAndRunServer(logger)
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
		logger.Infof(strings.LogGetSignalSuccess, signal.String())

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

func initAndRunServer(logger *logging.Logger) {

	// Получение конфигурации
	dbConfig := config.NewConfigDb()
	// Коннект к БД
	authDatabase, err := postgresql.NewAuthStorage(
		context.Background(),
		dbConfig,
		logger,
	)
	if err != nil {
		logger.Fatalln(strings.LogGetDatabaseError)
	}

	// Реализации интерфейсов для работы с БД
	tokenStorage := &postgresql.TokenDatabase{
		AuthDatabase: authDatabase,
	}
	userStorage := &postgresql.UserDatabase{
		AuthDatabase: authDatabase,
	}

	// Получение сервисов
	tokenService := &services.TokenService{
		TokenStorage: tokenStorage,
		UserStorage:  userStorage,
		Logger:       logger,
	}
	userService := &services.UserService{
		UserStorage:  userStorage,
		TokenService: tokenService,
		Logger:       logger,
	}

	// Композиция сервисов
	authService := &services.AuthService{
		TokenService: tokenService,
		UserService:  userService,
	}

	logger.Infoln(strings.LogInitRouters)
	router := mux.NewRouter()

	apiRouter := &routers.ApiRouter{
		Router:      router,
		AuthService: authService,
		Logger:      logger,
	}

	// Инициализация сервера
	serv := server.NewServer(
		apiRouter,
		authService,
		logger,
	)

	// Запуск сервера
	serv.Run()
}
