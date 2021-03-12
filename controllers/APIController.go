package controllers

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"net/http"
)

func APIController (r chi.Router)  {
	r.Get("/commands", func(w http.ResponseWriter, r *http.Request) {
		res, err := helpers.GetRedisClient().Get(context.Background(), "commandModules").Result()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(res))

	})


	r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		res, err := helpers.GetRedisClient().Get(context.Background(), "stats").Result()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(res))

	})
}

