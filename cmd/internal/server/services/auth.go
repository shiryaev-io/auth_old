package services

// Композиция сервисов
type AuthService struct {
	TokenService *TokenService
	UserService  *UserService
}
