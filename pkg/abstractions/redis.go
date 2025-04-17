package abstractions

import "github.com/redis/go-redis/v9"

type RedisClient interface {
	redis.UniversalClient
}
