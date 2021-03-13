package helpers

import (
	"github.com/utillybot/server/discord"
	"strconv"
)

func IsManageable(guild discord.PartialGuild) bool {
	permissions, _ := strconv.Atoi(guild.Permissions)
	if guild.Owner {
		return true
	}
	if permissions&0x00000008 == 1 {
		return true
	}

	return false
}
