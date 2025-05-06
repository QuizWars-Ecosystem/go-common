package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

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
	if cfg.Masters < 3 {
		cfg.Masters = 3
	}

	if cfg.Replicas >= cfg.Masters {
		cfg.Replicas = 1
	}

	totalNodes := cfg.Masters + cfg.Replicas*cfg.Masters

	exposedPorts := make([]string, totalNodes)
	for i := 0; i < totalNodes; i++ {
		exposedPorts[i] = fmt.Sprintf("%d/tcp", 7000+i)
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        cfg.Image,
			ExposedPorts: exposedPorts,
			WaitingFor:   wait.ForLog("Cluster state changed: ok").WithStartupTimeout(2 * time.Minute),
			Env: map[string]string{
				"IP":                        "0.0.0.0",
				"INITIAL_PORT":              "7000",
				"CLUSTER_ONLY":              "true",
				"STANDALONE":                "false",
				"MASTERS":                   fmt.Sprint(cfg.Masters),
				"REPLICAS":                  fmt.Sprint(cfg.Replicas),
				"REDIS_CLUSTER_CREATOR":     "yes",
				"REDIS_CLUSTER_REPLICAS":    fmt.Sprint(cfg.Replicas),
				"REDIS_CLUSTER_ANNOUNCE_IP": "host.docker.internal",
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				portBindings := nat.PortMap{}
				for i := 0; i < totalNodes; i++ {
					portBindings[nat.Port(exposedPorts[i])] = []nat.PortBinding{
						{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", 7000+i)},
					}
				}
				hc.PortBindings = portBindings
			},
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start redis cluster genericContainer: %w", err)
	}

	return genericContainer, nil
}
