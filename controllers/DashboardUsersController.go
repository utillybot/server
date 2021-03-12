package controllers

import (
	"github.com/go-chi/chi"
	"github.com/utillybot/server/middlewares"
	"io/ioutil"
	"net/http"
)

func DashboardUsersController(r chi.Router) {
	r.Use(middlewares.IsAuthenticated)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := middlewares.GetSessionData(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token := session.AccessToken

		req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v8/users/@me", nil)

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

		body, err := ioutil.ReadAll(result.Body)
		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(body)
	})
}