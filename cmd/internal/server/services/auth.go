package services

// Интерфейс, который должна реализовывать БД
type AuthStorage interface {
	// TODO: методы для работы сервиса
}

type AuthService struct {
	storage AuthStorage
}

func NewAuthService(storage AuthStorage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}
