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
func UnauthorizedError() *ApiError {
	return &ApiError{
		Status:     http.StatusUnauthorized,
		Err:        nil,
		Message:    "Пользователь не авторизован", // TODO: вынести сообщение в strings
		DevMessage: strings.Empty,
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
