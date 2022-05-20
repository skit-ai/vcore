package redis

import (
	"fmt"
	"github.com/mediocregopher/radix/v3"
	//"github.com/mediocregopher/radix.v2/redis"
	"os"
)

var (
	// Client - Redis client
	Client         *RadixRedisClient
	OutboundClient *RadixRedisClient
)

// RadixRedisClient - Redis client
type RadixRedisClient struct {
	*radix.Pool
}

// NewRadixRedisClient - Return a redis client
func NewRadixRedisClient() (*RadixRedisClient, error) {
	redisClient, err := radix.NewPool("tcp", fmt.Sprintf("%v:%v", os.Getenv("REDIS_REMOTE_HOST"), os.Getenv("REDIS_REMOTE_PORT")), 10)
	if err != nil {
		return nil, err
	}
	return &RadixRedisClient{redisClient}, nil
}

// NewRadixRedisClient - Return a redis client
func NewRadixRedisPool(size int, opts ...radix.PoolOpt) (*RadixRedisClient, error) {
	redisClient, err := radix.NewPool("tcp", fmt.Sprintf("%v:%v", os.Getenv("REDIS_REMOTE_HOST"), os.Getenv("REDIS_REMOTE_PORT")), size, opts...)
	if err != nil {
		return nil, err
	}
	return &RadixRedisClient{redisClient}, nil
}

// NewRadixRedisClient - Return a redis client
func NewRadixRedisClientUsingCustomHostPort(host, port string) (*RadixRedisClient, error) {
	redisClient, err := radix.NewPool("tcp", fmt.Sprintf("%v:%v", host, port), 10)
	if err != nil {
		return nil, err
	}
	return &RadixRedisClient{redisClient}, nil
}
