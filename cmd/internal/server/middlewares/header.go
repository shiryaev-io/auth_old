package middlewares

import "net/http"

// Устанавливает необходимые заголовки
func HeaderMiddleare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		
		handler.ServeHTTP(response, request)
	})
}