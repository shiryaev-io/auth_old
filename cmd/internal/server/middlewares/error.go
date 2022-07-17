package middlewares

import (
	"auth/cmd/internal/server/exceptions"
	"errors"
	"net/http"
)

// Функция для контроллеров, которая возвращает ошибку
type handler func(response http.ResponseWriter, request *http.Request) error

// Обрабатывает ошибки
func ErrorMiddleware(next handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		err := next(response, request)
		var apiError *exceptions.ApiError
		if err != nil {
			// Если ошибка err является кастомной ошибкой ApiError
			if errors.As(err, &apiError) {
				apiError = err.(*exceptions.ApiError)

				response.WriteHeader(apiError.Status)
				response.Write(apiError.Marshal())
			}

			// Если ошибка не кастомная, то возвращаем 500 статус
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("Непредвиденная ошибка")) // TODO: вынести сообщение в strings
		}
	}
}
