package redis

import (
  "os"
  "fmt"

  "github.com/go-redis/redis"
)

type GoRedisClient struct{
  *redis.Client
}

var (
  RedisClient *GoRedisClient
)

func NewGoRedisClient()(*GoRedisClient) {
  redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_REMOTE_HOST"), os.Getenv("REDIS_REMOTE_PORT")),
	})
  return &GoRedisClient{redisClient}
}
