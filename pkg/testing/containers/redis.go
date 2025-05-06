package containers

import (
	"context"
	"fmt"
	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"os/exec"
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

func NewRedisClusterContainersV2(ctx context.Context, cfg *config.RedisClusterConfig) ([]testcontainers.Container, []string, error) {
	var containers []testcontainers.Container
	var containerURLs []string

	for i := 0; i < cfg.ClusterSize; i++ {
		req := testcontainers.ContainerRequest{
			Image:        cfg.Image,
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379"),
			Env: map[string]string{
				"REDIS_CLUSTER": "yes",
			},
			Cmd: []string{
				"redis-server",
				"--port", "6379",
				"--cluster-enabled", "yes",
				"--cluster-config-file", fmt.Sprintf("nodes-%d.conf", i),
				"--cluster-node-timeout", "5000",
				"--appendonly", "yes",
			},
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

		if err != nil {
			return nil, nil, fmt.Errorf("failed to start redis container: %w", err)
		}

		host, err := container.Host(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get container host: %w", err)
		}
		port, err := container.MappedPort(ctx, "6379")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get container port: %w", err)
		}

		containerURLs = append(containerURLs, fmt.Sprintf("%s:%s", host, port.Port()))

		containers = append(containers, container)
	}

	for i := 0; i < len(containerURLs); i++ {
		for j := i + 1; j < len(containerURLs); j++ {
			cmd := fmt.Sprintf(
				"redis-cli --cluster add-node %s %s", containerURLs[j], containerURLs[i],
			)
			// Здесь запускаем команду для объединения Redis контейнеров в кластер
			_, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to add node to cluster: %w", err)
			}
		}
	}

	return containers, containerURLs, nil
}
