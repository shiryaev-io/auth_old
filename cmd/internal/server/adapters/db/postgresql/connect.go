package postgresql

import (
	"auth/cmd/internal/res/strings"
	"auth/cmd/internal/server/config"
	"auth/cmd/pkg/logging"
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
func connectDb(
	ctx context.Context,
	dbConfig *config.ConfigDb,
	logger *logging.Logger,
) (pool *pgxpool.Pool, err error) {

	dsn := fmt.Sprintf(
		dbUrl,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	retry.ForeverSleep(
		timeDelay,
		func(attempt int) error {

			logger.Infoln(strings.LogAttemptConnectDb, attempt)

			ctx, cancel := context.WithTimeout(ctx, timeDelay)
			defer cancel()

			logger.Infoln(strings.LogTryConnectDb)

			pool, err = pgxpool.Connect(ctx, dsn)
			if err != nil {
				logger.Fatalf(strings.LogFatalConnectDb, err)
				return err
			} else {
				logger.Infoln(strings.LogConnectSuccess)
			}

			return nil
		},
	)
	// TODO: возможно убрать/перенести в другое место метод Close()
	defer pool.Close()

	return pool, nil
}
