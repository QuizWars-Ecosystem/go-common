package containers

import (
	"context"
	"fmt"
	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
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

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        cfg.Image,
			ExposedPorts: []string{"6379-6384:6379-6384"},
			Env: map[string]string{
				"REDIS_NODES":               fmt.Sprint(cfg.Nodes),
				"REDIS_CLUSTER_REPLICAS":    fmt.Sprint(cfg.Replicas),
				"REDIS_CLUSTER_ANNOUNCE_IP": "host.docker.internal",
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("Cluster state changed: ok").WithStartupTimeout(time.Minute*5),
				wait.ForListeningPort("6379").WithStartupTimeout(time.Minute*5),
			),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start redis cluster container: %w", err)
	}

	return container, nil
}
