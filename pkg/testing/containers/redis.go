package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRedisContainer(ctx context.Context, cfg *config.RedisConfig) (*redis.RedisContainer, error) {
	return redis.Run(
		ctx,
		cfg.Image,
		testcontainers.WithHostPortAccess(cfg.Ports...),
	)
}

func NewRedisClusterContainers(ctx context.Context, cfg *config.RedisClusterConfig) (testcontainers.Container, error) {
	if cfg.Nodes < 3 {
		cfg.Nodes = 3
	}

	if cfg.Replicas >= cfg.Nodes {
		cfg.Replicas = 1
	}

	exposedPorts := make([]string, cfg.Nodes)
	for i := 0; i < cfg.Nodes; i++ {
		exposedPorts[i] = fmt.Sprintf("%d/tcp", 6379+i)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        cfg.Image,
			ExposedPorts: exposedPorts,
			Env: map[string]string{
				"REDIS_NODES":               fmt.Sprint(cfg.Nodes),
				"REDIS_CLUSTER_REPLICAS":    fmt.Sprint(cfg.Replicas),
				"REDIS_CLUSTER_ANNOUNCE_IP": "127.0.0.1",
				"ALLOW_EMPTY_PASSWORD":      "yes",
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("Ready to accept connections").WithStartupTimeout(2*time.Minute),
				wait.ForListeningPort("6379").WithStartupTimeout(2*time.Minute),
			),
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start redis cluster container: %w", err)
	}

	return container, nil
}
