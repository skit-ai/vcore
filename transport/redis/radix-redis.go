package redis

import (
	"fmt"
	"log"
	"os"

	"github.com/mediocregopher/radix.v2/redis"
)

type RadixRedisClient struct {
	*redis.Client
}

var (
	Client *RadixRedisClient
)

func NewRadixRedisClient() (*RadixRedisClient, error) {
	redisClient, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", os.Getenv("REDIS_REMOTE_HOST"), os.Getenv("REDIS_REMOTE_PORT")))
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	return &RadixRedisClient{redisClient}, nil
}
