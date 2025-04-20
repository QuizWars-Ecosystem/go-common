package config

import (
	"time"
)

type GatewayConfig struct {
	Local        bool          `mapstructure:"local"`
	LogLovel     string        `mapstructure:"log_lovel"`
	HTTPPort     int           `mapstructure:"http_port"`
	TCPPort      int           `mapstructure:"tcp_port"`
	GRPCPort     int           `mapstructure:"grpc_port"`
	WSPort       int           `mapstructure:"ws_port"`
	StartTimeout time.Duration `mapstructure:"start_timeout"`
	StopTimeout  time.Duration `mapstructure:"stop_timeout"`
	ConsulURL    string        `mapstructure:"consul_url"`
}

type ServiceConfig struct {
	Name         string        `mapstructure:"name"`
	Address      string        `mapstructure:"address"`
	Local        bool          `mapstructure:"local"`
	LogLovel     string        `mapstructure:"log_lovel"`
	GRPCPort     int           `mapstructure:"grpc_port"`
	StartTimeout time.Duration `mapstructure:"start_timeout"`
	StopTimeout  time.Duration `mapstructure:"stop_timeout"`
	ConsulURL    string        `mapstructure:"consul_url"`
}

type DefaultGatewayConfig struct {
	Local           bool          `env:"LOCAL" envDefault:"true"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"info"`
	HTTPPort        string        `env:"HTTP_PORT" envDefault:"8000"`
	TCPPort         string        `env:"TCP_PORT" envDefault:"8001"`
	GRPCPort        string        `env:"GRPC_PORT" envDefault:"8002"`
	WSPort          string        `env:"WS_PORT" envDefault:"8003"`
	StartTimeout    time.Duration `env:"START_TIMEOUT" envDefault:"15s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"15s"`
	ConsulURL       string        `env:"CONSUL_URL" envDefault:"http://127.0.0.1:8500"`
}

type DefaultServiceConfig struct {
	Name            string        `env:"NAME" envDefault:"service"`
	Address         string        `env:"ADDRESS" envDefault:"127.0.0.1"`
	Local           bool          `env:"LOCAL" envDefault:"true"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"info"`
	GRPCPort        int           `env:"GRPC_PORT" envDefault:"50000"`
	StartTimeout    time.Duration `env:"START_TIMEOUT" envDefault:"15s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"15s"`
	ConsulURL       string        `env:"CONSUL_URL" envDefault:"http://127.0.0.1:8500"`
}
