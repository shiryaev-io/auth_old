package controllers

import "net/http"

// Авторизация пользователя
func SigIn(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hello, world"))
}