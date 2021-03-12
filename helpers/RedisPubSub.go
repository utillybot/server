package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var guilds []string
var redisClient *redis.Client

func StartRedis() {
	url, _ := os.LookupEnv("REDIS_URL")
	opt, _ := redis.ParseURL(url)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       0,
	})

	defer redisClient.Close()

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}

	ctx := context.Background()
	// Subscribe to the Topic given
	topic := redisClient.Subscribe(ctx, "guilds", "guildCreate")
	// Get the Channel to use
	channel := topic.Channel()
	// Iterate any messages sent on the channel
	for msg := range channel {
		fmt.Println("Receiving message " + msg.Channel)
		if msg.Channel == "guilds" {
			payload := &guildsPayload{}

			if err := payload.UnmarshalBinary([]byte(msg.Payload)); err != nil {
				panic(err)
			}

			guilds = payload.Guilds
		} else if msg.Channel == "guildCreate" {
			payload := &guildCreatePayload{}

			if err := payload.UnmarshalBinary([]byte(msg.Payload)); err != nil {
				panic(err)
			}

			guilds = append(guilds, payload.Guild)
		}

	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func GetGuilds() []string {
	return guilds
}

type guildsPayload struct {
	Guilds []string `json:"guilds"`
}

type guildCreatePayload struct {
	Guild string `json:"guild"`
}

func (u *guildsPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *guildsPayload) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	return nil
}

func (u *guildCreatePayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *guildCreatePayload) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	return nil
}
