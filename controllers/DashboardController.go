package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/utillybot/server/helpers"
	"github.com/utillybot/server/middlewares"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func DashboardController(r chi.Router) {
	scopes := []string{"identify", "guilds"}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	r.Route("/users", DashboardUsersController)
	r.Route("/guilds", DashboardGuildsController)

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		middlewares.DestroySession(r.Context())

		if err := middlewares.GetSession(r.Context()).Save(r, w); err != nil {
			http.StatusText(http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/",302)
		}
	})

	r.Get("/authorize", func(w http.ResponseWriter, r *http.Request) {
		requestUrl, _ := url.Parse("https://discord.com/api/oauth2/authorize")

		query := requestUrl.Query()
		query.Set("response_type", "code")
		query.Set("client_id", clientID)
		query.Set("scope", strings.Join(scopes, " "))
		query.Set("redirect_uri", "http://"+r.Host+"/api/dashboard/callback")
		requestUrl.RawQuery = query.Encode()

		http.Redirect(w, r, requestUrl.String(), 302)
	})

	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("error") != "" {
			http.Redirect(w, r, "/api/dashboard/done", 302)
			return
		}

		code := query.Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("No access code provided."))
			return
		}

		obj := url.Values{}
		obj.Set("client_id", clientID)
		obj.Set("client_secret", clientSecret)
		obj.Set("grant_type", "authorization_code")
		obj.Set("code", code)
		obj.Set("redirect_uri", "http://"+r.Host+"/api/dashboard/callback")
		obj.Set("scope", strings.Join(scopes, " "))

		response, err := http.PostForm("https://discord.com/api/oauth2/token", obj)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var result helpers.SessionData

		if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if result.Scope != strings.Join(scopes, " ") {
			http.Redirect(w, r, "/api/dashboard/authorize", 302)
		} else {
			session := middlewares.GetSession(r.Context())

			session.Values["AccessToken"] = result.AccessToken
			session.Values["ExpiresIn"] = result.ExpiresIn
			session.Values["RefreshToken"] = result.RefreshToken
			session.Values["Scope"] = result.Scope
			session.Values["TokenType"] = result.TokenType

			session.Options.MaxAge = result.ExpiresIn

			if err = session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard/done", 302)
		}
	})

}
