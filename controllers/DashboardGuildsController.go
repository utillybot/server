package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/discord"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/middlewares"
	"net/http"
)

func DashboardGuildsController() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.IsAuthenticated)
	r.Route("/{id}", DashboardGuildController)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		token, err := middlewares.GetAccessToken(r.Context())
		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		guilds, err := discord.GetGuilds(token)
		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		var present []MappedGuild
		var notPresent []MappedGuild

		for _, guild := range guilds {
			if helpers.IsManageable(guild) {
				mappedGuild := MappedGuild{
					Name: guild.Name,
					Icon: guild.Icon,
					Id:   guild.Id,
				}

				if middlewares.GuildExists(mappedGuild.Id) {
					present = append(present, mappedGuild)
				} else {
					notPresent = append(notPresent, mappedGuild)
				}
			}
		}

		getGuildsResponse := GetGuildsResponse{
			Present:    present,
			NotPresent: notPresent,
		}

		res, err := json.Marshal(getGuildsResponse)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(res)
	})

	return r
}

type MappedGuild struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Id   string `json:"id"`
}

type GetGuildsResponse struct {
	Present    []MappedGuild `json:"present"`
	NotPresent []MappedGuild `json:"notPresent"`
}
