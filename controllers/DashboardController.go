package controllers

import (
	"github.com/go-chi/chi"
	"github.com/utillybot/server/discord"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/middlewares"
	"net/http"
	"os"
	"strings"
)

func DashboardController() chi.Router {
	r := chi.NewRouter()

	scopes := []string{"identify", "guilds"}

	oauth2 := discord.OAuth2{
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scope:        scopes,
	}

	r.Mount("/users", DashboardUsersController())
	r.Mount("/guilds", DashboardGuildsController())

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		middlewares.DestroySession(r.Context())

		if err := middlewares.GetSession(r.Context()).Save(r, w); err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	})

	r.Get("/authorize", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, oauth2.GenerateAuthURL(discord.AuthURLOptions{
			RedirectUri: "http://" + r.Host + "/api/dashboard/callback",
		}), 302)
	})

	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("error") != "" {
			http.Redirect(w, r, "/api/dashboard/done", 302)
			return
		}

		code := query.Get("code")
		if code == "" {
			http.Error(w, "No access code provided.", http.StatusBadRequest)
			return
		}

		result, err := oauth2.TokenRequest(discord.TokenRequestOptions{
			Code:        code,
			RedirectUri: "http://" + r.Host + "/api/dashboard/callback",
		})

		if err != nil {
			helpers.HttpError(w, http.StatusInternalServerError)
			return
		}

		if result.Scope != strings.Join(scopes, " ") {
			http.Redirect(w, r, "/api/dashboard/authorize", 302)
		} else {
			session := middlewares.GetSession(r.Context())
			session.Values["Tokens"] = result
			session.Options.MaxAge = result.ExpiresIn

			if err = session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard/done", 302)
		}
	})

	return r
}
