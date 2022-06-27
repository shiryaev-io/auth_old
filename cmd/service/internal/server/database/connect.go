package database

import (
	"auth/cmd/service/pkg/logging"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tinrab/retry"
)

const (
	dbUrl     = "postgresql://%s:%s@%s:%s/%s"
	timeDelay = 5 * time.Second
)

// Подключение к БД
func ConnectDb(
	ctx context.Context,
	dbConfig *dbConfig,
	logger *logging.Logger,
) (pool *pgxpool.Pool, err error) {

	dsn := fmt.Sprintf(
		dbUrl,
		dbConfig.user,
		dbConfig.password,
		dbConfig.host,
		dbConfig.port,
		dbConfig.name,
	)

	retry.ForeverSleep(timeDelay, func(attempt int) error {

		logger.Infoln(logAttemptConnectDb, attempt)

		ctx, cancel := context.WithTimeout(ctx, timeDelay)
		defer cancel()

		logger.Infoln(logTryConnectDb)

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			logger.Fatalf(logFatalConnectDb, err)
			return err
		} else {
			logger.Infoln(logConnectSuccess)
		}

		return nil
	})

	return pool, nil
}
