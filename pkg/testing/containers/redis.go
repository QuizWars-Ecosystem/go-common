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

func NewRedisClusterContainers(ctx context.Context, cfg *config.RedisClusterConfig) ([]*redis.RedisContainer, error) {
	var containers []*redis.RedisContainer

	for i := 0; i < cfg.ClusterSize; i++ {
		container, err := redis.Run(
			ctx,
			cfg.Image,
			testcontainers.WithStartupCommand(
				testcontainers.NewRawCommand([]string{
					"redis-server",
					"--port", "6379",
					"--cluster-enabled", "yes",
					"--cluster-config-file", fmt.Sprintf("nodes-%d.conf", i),
					"--cluster-node-timeout", "5000",
					"--appendonly", "yes",
				}),
			),
			testcontainers.WithExposedPorts("6379", "16379"),
			testcontainers.WithWaitStrategy(wait.ForListeningPort("6379")),
		)

		if err != nil {
			return nil, fmt.Errorf("failed to start redis container: %w", err)
		}

		containers = append(containers, container)
	}

	return containers, nil
}
