package postgresql

import (
	"auth/cmd/internal/server/config"
	"auth/cmd/pkg/logging"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Структура базы данных
type authStorage struct {
	pool   *pgxpool.Pool
	logger *logging.Logger
}

func NewAuthStorage(
	ctx context.Context,
	config *config.ConfigDb,
	logger *logging.Logger,
) (*authStorage, error) {

	pool, err := connectDb(ctx, config, logger)
	if err != nil {
		return nil, err
	}

	return &authStorage{
		pool:   pool,
		logger: logger,
	}, nil
}
