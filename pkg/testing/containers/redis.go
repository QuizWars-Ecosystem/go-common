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

	exposedPorts := make([]string, cfg.Masters)
	for i := 0; i < cfg.Masters; i++ {
		exposedPorts[i] = fmt.Sprintf("%d/tcp", 7000+i)
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        cfg.Image,
			ExposedPorts: []string{"7000/tcp", "7001/tcp", "7002/tcp", "7003/tcp", "7004/tcp", "7005/tcp"},
			WaitingFor:   wait.ForLog("Cluster state changed: ok").WithStartupTimeout(2 * time.Minute),
			Env: map[string]string{
				"IP":                        "0.0.0.0",
				"INITIAL_PORT":              "7000",
				"CLUSTER_ONLY":              "true",
				"STANDALONE":                "false",
				"MASTERS":                   fmt.Sprint(cfg.Masters),
				"REPLICAS":                  fmt.Sprint(cfg.Replicas),
				"REDIS_CLUSTER_CREATOR":     "yes",
				"REDIS_CLUSTER_REPLICAS":    "1",
				"REDIS_CLUSTER_ANNOUNCE_IP": "host.docker.internal",
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"7000/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7000"}},
					"7001/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7001"}},
					"7002/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7002"}},
					"7003/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7003"}},
					"7004/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7004"}},
					"7005/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7005"}},
				}
			},
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start redis cluster genericContainer: %w", err)
	}

	return genericContainer, nil
}
