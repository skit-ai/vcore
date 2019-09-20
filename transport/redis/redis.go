package redis

import (
	"fmt"
	"os"

	"github.com/mediocregopher/radix.v2/redis"
)

var (
	// Client - Redis client
	Client *RadixRedisClient
)

// RadixRedisClient - Redis client
type RadixRedisClient struct {
	*redis.Client
}

// NewRadixRedisClient - Return a redis client
func NewRadixRedisClient() (*RadixRedisClient, error) {
	redisClient, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", os.Getenv("REDIS_REMOTE_HOST"), os.Getenv("REDIS_REMOTE_PORT")))
	if err != nil {
		return nil, err
	}
	return &RadixRedisClient{redisClient}, nil
}
