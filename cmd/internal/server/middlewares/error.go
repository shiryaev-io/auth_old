package middlewares

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"errors"
	"net/http"
)

// Функция handler, которая возвращает ошибку
type ErrorHandlerFunc func(response http.ResponseWriter, request *http.Request) error

// Обрабатывает ошибки
func ErrorMiddleware(next ErrorHandlerFunc) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		err := next(response, request)
		handleError(response, err)
	})
}

func handleError(response http.ResponseWriter, err error) {
	var apiError *exceptions.ApiError
	// Если ошибка err является кастомной ошибкой ApiError
	if errors.As(err, &apiError) {
		apiError = err.(*exceptions.ApiError)

		response.WriteHeader(apiError.Status)
		response.Write(apiError.Marshal())
	}

	// Если ошибка не кастомная, то возвращаем 500 статус
	status := http.StatusInternalServerError
	apiError = &exceptions.ApiError{
		Status:     status,
		Err:        err,
		Message:    strings.MessageUnforeseenError,
		DevMessage: err.Error(),
	}
	response.WriteHeader(status)
	response.Write(apiError.Marshal())
}
