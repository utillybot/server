package redisClient

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func CreateRedisClient() {
	url, _ := os.LookupEnv("REDIS_URL")
	opt, _ := redis.ParseURL(url)
	client = redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := client.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
}


func GetRedisClient() *redis.Client {
	return client
}
