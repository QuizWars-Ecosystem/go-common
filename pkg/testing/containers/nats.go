package containers

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
)

func NewNATSContainer(ctx context.Context, cfg *config.NatsConfig) (*nats.NATSContainer, error) {
	return nats.Run(
		ctx,
		cfg.Image,
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("Server is ready"),
				wait.ForListeningPort("4222/tcp"),
			).WithDeadline(10*time.Second),
		),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.PortBindings = nat.PortMap{
				"4222/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "4222"}},
			}
		}),
	)
}
