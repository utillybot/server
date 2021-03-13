package controllers

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/redisClient"
	"net/http"
)

func APIController(r chi.Router) {
	r.Mount("/dashboard", DashboardController())

	r.Get("/commands", func(w http.ResponseWriter, r *http.Request) {
		res, err := redisClient.GetRedisClient().Get(context.Background(), "commandModules").Result()

		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(res))

	})
	r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		res, err := redisClient.GetRedisClient().Get(context.Background(), "stats").Result()

		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(res))

	})
}
