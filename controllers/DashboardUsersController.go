package controllers

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/middlewares"
	"net/http"
)

func DashboardUsersController() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.IsAuthenticated)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		token, err := middlewares.GetAccessToken(r.Context())
		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}
		user, err := middlewares.GetCurrentUser(r.Context())

		if user == nil || err != nil {
			discordClient, err := discordgo.New("Bearer " + token)
			if err != nil {
				helpers.HttpError(w, http.StatusInternalServerError)
				return
			}
			user, err = discordClient.User("@me")
			if err != nil {
				helpers.HttpError(w, http.StatusInternalServerError)
				return
			}

			session := middlewares.GetSession(r.Context())
			session.Values["User"] = user
			err = session.Save(r, w)
			if err != nil {
				helpers.HttpError(w, http.StatusInternalServerError)
				return
			}
		}

		response, err := json.Marshal(user)
		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(response)
	})

	return r
}
