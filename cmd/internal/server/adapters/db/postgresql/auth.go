package postgresql

import (
	"auth/cmd/internal/server/config"
	"auth/cmd/pkg/logging"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Структура базы данных
type authDatabase struct {
	pool   *pgxpool.Pool
	logger *logging.Logger
}

func NewAuthStorage(
	ctx context.Context,
	config *config.ConfigDb,
	logger *logging.Logger,
) (*authDatabase, error) {

	pool, err := connectDb(ctx, config, logger)
	if err != nil {
		return nil, err
	}

	return &authDatabase{
		pool:   pool,
		logger: logger,
	}, nil
}
