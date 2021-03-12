package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/middlewares"
	"net/http"
)

func DashboardGuildsController(r chi.Router) {
	r.Use(middlewares.IsAuthenticated)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := middlewares.GetSessionData(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token := session.AccessToken
		req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v8/users/@me/guilds", nil)

		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer " + token)

		result, err := http.DefaultClient.Do(req)

		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var guilds []helpers.PartialGuild

		err = json.NewDecoder(result.Body).Decode(&guilds)
		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

	r.Route("/{id}", func(r chi.Router) {
		r.Use(middlewares.ValidateGuild)

		r.Get("/",  func(w http.ResponseWriter, r *http.Request) {
			guild := middlewares.GetGuild(r.Context())

			res, err := json.Marshal(guild)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			_, _ = w.Write(res)
		})
	})

}

type MappedGuild struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Id string `json:"id"`
}

type GetGuildsResponse struct {
	Present    []MappedGuild `json:"present"`
	NotPresent []MappedGuild `json:"notPresent"`
}

