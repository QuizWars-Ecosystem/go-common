package containers

import (
	"context"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
)

func NewNATSContainer(ctx context.Context, cfg *config.NatsConfig) (*nats.NATSContainer, error) {
	return nats.Run(
		ctx,
		cfg.Image,
		testcontainers.WithExposedPorts("4222"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("4222")),
	)
}
