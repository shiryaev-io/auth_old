package exceptions

import (
	"auth/cmd/internal/res/strings"
	"encoding/json"
	"net/http"
)

// Кастомная ошибка запросов
type ApiError struct {
	Status     int
	Err        error
	Message    string
	DevMessage string
}

// Реализация интрейфейса Error{}
func (apiError *ApiError) Error() string {
	return apiError.Message
}

// Возвращает массив байтов для передачи на клиент
func (apiError *ApiError) Marshal() []byte {
	marshal, err := json.Marshal(apiError)
	if err != nil {
		return nil
	}
	return marshal
}

// Возвращает ошибку, что пользователь не авторизован
func UnauthorizedError(err error) *ApiError {
	return &ApiError{
		Status:     http.StatusUnauthorized,
		Err:        err,
		Message:    strings.ErrorUserUnauthorized,
		DevMessage: err.Error(),
	}
}

// Возвращает ошибку, если пользователь ввел неккоретные данные, 
// не прошел валидацию и т.д.
func BadRequest(message string, err error) *ApiError {
	return &ApiError{
		Status:     http.StatusBadRequest,
		Err:        err,
		Message:    message,
		DevMessage: err.Error(),
	}
}

// Ошибка сервера: не сгенерировался токен и т.д.
func ServerError(message string, err error) *ApiError {
	return &ApiError{
		Status:     http.StatusInternalServerError,
		Err:        err,
		Message:    message,
		DevMessage: err.Error(),
	}
}
