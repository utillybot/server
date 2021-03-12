package helpers

import "strconv"

func IsManageable(guild PartialGuild) bool {
	permissions, _ := strconv.Atoi(guild.Permissions)
	if guild.Owner {
		return true
	}
	if permissions & 0x00000008 == 1 {
		return true
	}

	return false
}