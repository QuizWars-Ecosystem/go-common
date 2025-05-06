package config

type PostgresConfig struct {
	Name     string
	Image    string
	Username string
	Password string
	DBName   string
	Host     string
	Ports    []int
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Name:     "postgres",
		Image:    "postgres:17.2-alpine",
		Username: "user",
		Password: "pass",
		DBName:   "db",
		Host:     "localhost",
		Ports:    []int{5432},
	}
}

type RedisConfig struct {
	Name  string
	Image string
	Host  string
	Ports []int
}

func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Name:  "redis",
		Image: "redis:7.4-alpine",
		Host:  "localhost",
		Ports: []int{6379},
	}
}

type RedisClusterConfig struct {
	RedisConfig
	ClusterSize int
}

func DefaultRedisClusterConfig() RedisClusterConfig {
	return RedisClusterConfig{
		RedisConfig: DefaultRedisConfig(),
		ClusterSize: 3,
	}
}

type ConsulConfig struct {
	Name  string
	Image string
	Host  string
	Ports []int
}

func DefaultConsulConfig() ConsulConfig {
	return ConsulConfig{
		Name:  "consul",
		Image: "hashicorp/consul:1.20",
		Host:  "localhost",
		Ports: []int{8500, 8600},
	}
}

type KafkaConfig struct {
	Name      string
	Image     string
	ClusterID string
	Ports     []int
}

func DefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Name:      "kafka",
		Image:     "confluentinc/confluent-local:7.9.0",
		ClusterID: "kafka-test",
		Ports:     []int{9092, 9093, 9094, 9095},
	}
}

type NatsConfig struct {
	Name      string
	Image     string
	ClusterID string
	Command   string
	Ports     []int
}

func DefaultNatsConfig() NatsConfig {
	return NatsConfig{
		Name:      "nats",
		Image:     "nats:2-alpine3.21",
		ClusterID: "nats-test",
		Command:   "command: -js -m 8222",
		Ports:     []int{4222, 8222},
	}
}
