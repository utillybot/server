package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/middlewares"
	"net/http"
)

func DashboardGuildController(r chi.Router) {
	r.Use(middlewares.ValidateGuild)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		guild := middlewares.GetCurrentGuild(r.Context())

		res, err := json.Marshal(guild)

		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(res)
	})
}
