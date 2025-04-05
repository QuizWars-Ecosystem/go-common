package containers

import (
	"context"

	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

func NewKafkaContainer(ctx context.Context, cfg *config.KafkaConfig) (*kafka.KafkaContainer, error) {
	return kafka.Run(
		ctx,
		cfg.Image,
		kafka.WithClusterID(cfg.ClusterID),
		testcontainers.WithHostPortAccess(cfg.Ports...),
	)
}
