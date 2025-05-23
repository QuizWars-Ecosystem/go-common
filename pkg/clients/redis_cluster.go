package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type RedisClusterOptions struct {
	*redis.ClusterOptions
	traceEnabled bool
	provider     trace.TracerProvider
}

var defaultClusterOptions = &redis.ClusterOptions{
	Addrs:           []string{"127.0.0.1:7000", "127.0.0.1:7001"},
	MaxRedirects:    3,
	ReadOnly:        true,
	RouteByLatency:  true,
	RouteRandomly:   false,
	PoolSize:        100,
	MinIdleConns:    10,
	DialTimeout:     10 * time.Second,
	ReadTimeout:     time.Second,
	WriteTimeout:    time.Second,
	MaxRetries:      5,
	MinRetryBackoff: 10 * time.Millisecond,
	MaxRetryBackoff: time.Second,
}

func NewRedisClusterOptions(addrs []string) *RedisClusterOptions {
	opts := defaultClusterOptions
	opts.Addrs = addrs
	return &RedisClusterOptions{
		ClusterOptions: opts,
	}
}

func NewRedisClusterClient(opts *RedisClusterOptions) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(opts.ClusterOptions)

	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var lastErr, err error
	for i := 0; i < opts.MaxRetries; i++ {
		if err = client.Ping(pingCtx).Err(); err == nil {
			break
		}

		lastErr = err
		time.Sleep(time.Second * time.Duration(i+1))
	}
	if lastErr != nil {
		return nil, fmt.Errorf("redis cluster ping failed: %w", lastErr)
	}

	if opts.traceEnabled {
		if err = redisotel.InstrumentTracing(client, redisotel.WithTracerProvider(opts.provider)); err != nil {
			return nil, fmt.Errorf("redis otel init failed: %w", err)
		}
	}

	return client, nil
}

func (o *RedisClusterOptions) WithAddrs(addrs []string) *RedisClusterOptions {
	o.Addrs = addrs
	return o
}

func (o *RedisClusterOptions) WithMaxRedirects(maxRedirects int) *RedisClusterOptions {
	o.MaxRedirects = maxRedirects
	return o
}

func (o *RedisClusterOptions) WithMaxRetries(maxRetries int) *RedisClusterOptions {
	o.MaxRetries = maxRetries
	return o
}

func (o *RedisClusterOptions) WithReadOnlyFlag(readOnly bool) *RedisClusterOptions {
	o.ReadOnly = readOnly
	return o
}

func (o *RedisClusterOptions) WithRouterByLatency(routerByLatency bool) *RedisClusterOptions {
	o.RouteByLatency = routerByLatency
	o.RouteRandomly = !routerByLatency
	return o
}

func (o *RedisClusterOptions) WithRouteRandomly(randomly bool) *RedisClusterOptions {
	o.RouteRandomly = randomly
	o.RouteByLatency = !randomly
	return o
}

func (o *RedisClusterOptions) WithPoolSize(poolSize int) *RedisClusterOptions {
	o.PoolSize = poolSize
	return o
}

func (o *RedisClusterOptions) WithDialTimeout(dialTimeout time.Duration) *RedisClusterOptions {
	o.DialTimeout = dialTimeout
	return o
}

func (o *RedisClusterOptions) WithMinIdleConns(minIdleConns int) *RedisClusterOptions {
	o.MinIdleConns = minIdleConns
	return o
}

func (o *RedisClusterOptions) WithReadTimeout(timeout time.Duration) *RedisClusterOptions {
	o.ReadTimeout = timeout
	return o
}

func (o *RedisClusterOptions) WithWriteTimeout(timeout time.Duration) *RedisClusterOptions {
	o.WriteTimeout = timeout
	return o
}

func (o *RedisClusterOptions) WithPoolTimeout(timeout time.Duration) *RedisClusterOptions {
	o.PoolTimeout = timeout
	return o
}

func (o *RedisClusterOptions) WithBackoffTimeouts(min, max time.Duration) *RedisClusterOptions {
	o.MinRetryBackoff = min
	o.MaxRetryBackoff = max
	return o
}

func (o *RedisClusterOptions) WithTraceProvider(provider trace.TracerProvider) *RedisClusterOptions {
	o.traceEnabled = true
	o.provider = provider
	return o
}
