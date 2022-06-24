package controllers

import "net/http"

// Авторизация пользователя
func SigIn(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello, world"))
}