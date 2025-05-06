package containers

import (
	"context"
	"fmt"
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
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: cfg.Image,
			ExposedPorts: []string{
				"7000-7005:7000-7005", // Ports Redis Cluster
			},
			Env: map[string]string{
				"INITIAL_PORT":              fmt.Sprint(cfg.InitialPort),
				"MASTERS":                   fmt.Sprint(cfg.MasterNodes),     // Count of master nodes
				"REPLICAS":                  fmt.Sprint(cfg.SlavesPerMaster), // Count of replica nod on each master node
				"IP":                        "0.0.0.0",
				"REDIS_CLUSTER_ANNOUNCE_IP": "host.docker.internal", // For access from host
			},
			WaitingFor: wait.ForLog("Cluster state changed: ok"),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start redis cluster container: %w", err)
	}

	return container, nil
}
