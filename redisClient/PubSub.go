package redisClient

import (
	"context"
	"encoding/json"
	"fmt"
)

var guildIds []string

func StartListening() {
	ctx := context.Background()
	// Subscribe to the Topic given
	topic := client.Subscribe(ctx, "guilds", "guildCreate")
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

			guildIds = payload.Guilds
		} else if msg.Channel == "guildCreate" {
			payload := &guildCreatePayload{}

			if err := payload.UnmarshalBinary([]byte(msg.Payload)); err != nil {
				panic(err)
			}

			guildIds = append(guildIds, payload.Guild)
		}
	}
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

func GetGuilds() []string {
	return guildIds
}
