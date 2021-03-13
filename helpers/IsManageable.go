package helpers

import (
	"github.com/bwmarrin/discordgo"
)

func IsManageable(guild *discordgo.UserGuild) bool {
	if guild.Owner {
		return true
	}
	if guild.Permissions&0x00000008 != 0 {
		return true
	}

	return false
}
