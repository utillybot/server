package middlewares

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/discord"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/redisClient"
	"net/http"
)

const contextKeyGuild = helpers.ContextKey("guild")

func ValidateGuild(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAccessToken(r.Context())
		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		guilds, err := discord.GetGuilds(token)
		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")

		var foundGuild discord.PartialGuild
		for _, v := range guilds {
			if v.Id == id {
				foundGuild = v
			}
		}

		if foundGuild.Id == "" {
			http.Error(w, "The provided guild id was not one of the user's guilds.", http.StatusNotFound)
			return
		}

		if !helpers.IsManageable(foundGuild) {
			http.Error(w, "You do not have the permission to get this guild.", http.StatusForbidden)
			return
		}

		if !GuildExists(foundGuild.Id) {
			http.Error(w, "The bot is not in the guild", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyGuild, &foundGuild)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCurrentGuild(ctx context.Context) *discord.PartialGuild {
	guild := ctx.Value(contextKeyGuild).(*discord.PartialGuild)
	return guild
}

func GuildExists(id string) bool {
	for _, v := range redisClient.GetGuilds() {
		if v == id {
			return true
		}
	}
	return false
}
