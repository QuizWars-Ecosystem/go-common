package containers

import (
	"context"
	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func NewPostgresContainer(ctx context.Context, cfg *config.PostgresConfig) (*postgres.PostgresContainer, error) {
	return postgres.Run(
		ctx,
		cfg.Image,
		postgres.WithDatabase(cfg.DBName),
		postgres.WithUsername(cfg.Username),
		postgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*5),
		),
		testcontainers.WithHostPortAccess(cfg.Ports...),
	)
}
