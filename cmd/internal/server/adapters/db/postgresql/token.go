package postgresql

import (
	"auth/cmd/internal/server/adapters/db/postgresql/queries"
	"auth/cmd/internal/server/models/db"
	"auth/cmd/internal/server/models/dto"
	"context"
	"errors"
)

// Сткруатруа для реелизации интерфейса TokenStorage
type TokenDatabase struct {
	AuthDatabase *authDatabase
}

// Получение модели токена из БД по userId
func (storage *TokenDatabase) FindByUserId(userId int) (*db.Token, error) {
	token := &db.Token{}
	query := queries.QuerySelectTokenByUserId
	err := storage.AuthDatabase.
		pool.
		QueryRow(context.Background(), query, userId).
		Scan(
			&token.Id,
			&token.UserId,
			&token.Value,
		)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Получение модели токена из БД
func (storage *TokenDatabase) FindToken(refreshToken string) (*db.Token, error) {
	token := &db.Token{}
	query := queries.QuerySelectTokenByRefreshToken
	err := storage.AuthDatabase.
		pool.
		QueryRow(context.Background(), query, refreshToken).
		Scan(
			&token.Id,
			&token.UserId,
			&token.Value,
		)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Сохранение токена в БД
func (storage *TokenDatabase) SaveToken(tokenDto *dto.Token) error {
	query := queries.QueryUpdateToken
	result, err := storage.AuthDatabase.
		pool.
		Exec(context.Background(), query, tokenDto.Value, tokenDto.Id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		// TODO: вынести с строки
		return errors.New("Ни одна строка не была обновлена")
	}
	return nil
}


// Создание новой записи в БД
func (storage *TokenDatabase) CreateToken(userId int, refreshToken string) (*db.Token, error) {
	token := &db.Token{
		UserId: userId,
		Value:  refreshToken,
	}
	query := queries.QueryInsertToken
	err := storage.AuthDatabase.
		pool.
		QueryRow(context.Background(), query, userId, refreshToken).
		Scan(&token.Id)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Удаляет токен из БД
func (storage *TokenDatabase) RemoveToken(refreshToken string) (*db.Token, error) {
	token := &db.Token{}
	query := queries.QueryDeleteToken
	err := storage.AuthDatabase.
		pool.
		QueryRow(context.Background(), query, refreshToken).
		Scan(
			&token.Id,
			&token.UserId,
			&token.Value,
		)
	if err != nil {
		return nil, err
	}
	return token, nil
}
