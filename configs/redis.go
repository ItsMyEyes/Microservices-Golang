package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func init() {
	connect()
}

var ClientRedis *redis.Client

// create connection to redis
func connect() *redis.Client {
	ClientRedis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	pong, err := ClientRedis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, "connected to rds")
	return ClientRedis
}

func SetKey(key string, i interface{}, ttl time.Duration) *redis.StringCmd {
	if ttl == 0 {
		ClientRedis.Set(context.Background(), key, i, 0)
	} else {
		ClientRedis.Set(context.Background(), key, i, ttl)
	}

	return ClientRedis.Get(context.Background(), key)
}

func GetKey(key string) *redis.StringCmd {
	return ClientRedis.Get(context.Background(), key)
}

func RemoveKey(key string) bool {
	ClientRedis.Del(context.Background(), key)
	return true
}
