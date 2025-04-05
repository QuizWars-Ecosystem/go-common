package containers

import (
	"context"
	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/consul"
)

func NewConsulContainer(ctx context.Context, cfg *config.ConsulConfig) (*consul.ConsulContainer, error) {
	return consul.Run(
		ctx,
		cfg.Image,
		testcontainers.WithStartupCommand(
			testcontainers.NewRawCommand([]string{"agent", "-dev", "-client=0.0.0.0"}),
		),
		testcontainers.WithHostPortAccess(cfg.Ports...),
	)
}
