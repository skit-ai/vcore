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

//TODO: Redis reconnection
func NewRadixRedisClient() *RadixRedisClient {
	redisClient, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", os.Getenv("REDIS_REMOTE_HOST"), os.Getenv("REDIS_REMOTE_PORT")))
	if err != nil {
		log.Fatalf("%s", err)
	}
	return &RadixRedisClient{redisClient}
}
