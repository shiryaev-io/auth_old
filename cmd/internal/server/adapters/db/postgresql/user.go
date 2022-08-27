package postgresql

import (
	"auth/cmd/internal/server/adapters/db/postgresql/queries"
	"auth/cmd/internal/server/models/db"
	"context"
)

// Структура для реализации интрефейса UserStorage
type UserDatabase struct {
	AuthDatabase *authDatabase
}

// Получение данных пользователя по Email
func (storage *UserDatabase) FindByEmail(email string) (*db.User, error) {
	user := &db.User{}
	// TODO: возможно с помощью fmt нужно сформировать строку, где подставить данные вместо `?`
	query := queries.QuerySelectUserByEmail
	err := storage.AuthDatabase.
		pool.
		QueryRow(context.Background(), query, email).
		Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.IsActivated,
		)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Получение данных пользователя по ID
func (storage *UserDatabase) FindById(userId int) (*db.User, error) {
	user := &db.User{}
	query := queries.QuerySelectUserById
	err := storage.AuthDatabase.
		pool.
		QueryRow(context.Background(), query, userId).
		Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.IsActivated,
		)
	if err != nil {
		return nil, err
	}
	return user, nil
}
