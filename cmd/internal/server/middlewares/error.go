package middlewares

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/exceptions"
	"auth/cmd/internal/server/models/responses"
	"errors"
	"net/http"
)

// Функция handler, которая возвращает ошибку
type ErrorHandlerFunc func(response http.ResponseWriter, request *http.Request) (*responses.Common, error)

// Обрабатывает ошибки
func ErrorMiddleware(next ErrorHandlerFunc) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		commonResponse, err := next(response, request)
		if err != nil {
			handleError(response, err)
			return
		}
		response.WriteHeader(commonResponse.Status)
		response.Write(commonResponse.Body)
	})
}

func handleError(response http.ResponseWriter, err error) {
	var apiError *exceptions.ApiError
	// Если ошибка err является кастомной ошибкой ApiError
	if errors.As(err, &apiError) {
		apiError = err.(*exceptions.ApiError)

		response.WriteHeader(apiError.Status)
		response.Write(apiError.Marshal())
		return
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
