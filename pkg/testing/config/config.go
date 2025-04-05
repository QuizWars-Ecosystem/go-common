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
