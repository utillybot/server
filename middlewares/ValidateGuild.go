package middlewares

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"net/http"
)

const contextKeyGuild = helpers.ContextKey("guild")

func ValidateGuild(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSessionData(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		token := session.AccessToken

		req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v8/users/@me/guilds", nil)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		result, err := http.DefaultClient.Do(req)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var guilds []helpers.PartialGuild

		err = json.NewDecoder(result.Body).Decode(&guilds)
		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")

		var foundGuild helpers.PartialGuild
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

func GetGuild (ctx context.Context) *helpers.PartialGuild {
	guild := ctx.Value(contextKeyGuild).(*helpers.PartialGuild)
	return guild
}

func GuildExists(id string) bool {
	for _, v := range helpers.GetGuilds() {
		if v == id {
			return true
		}
	}
	return false
}