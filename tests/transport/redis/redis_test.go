package tests

import (
	"testing"

	"github.com/skit-ai/vcore/transport/redis"
)

func TestNewRadixRedisClient(t *testing.T) {

	if _, err := redis.NewRadixRedisClient(); err != nil {
		t.Error(err)
	}
}
