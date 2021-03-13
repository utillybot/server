package middlewares

import (
	"context"
	"errors"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	"github.com/utillybot/server/discord"
	"github.com/utillybot/server/helpers"
	"net/http"
)

const contextKeySession = helpers.ContextKey("session")

func Sessions(store *pgstore.PGStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			ctx := context.WithValue(r.Context(), contextKeySession, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSession(ctx context.Context) *sessions.Session {
	session := ctx.Value(contextKeySession).(*sessions.Session)
	return session
}

func GetAccessToken(ctx context.Context) (string, error) {
	session := GetSession(ctx)

	tokens, ok := session.Values["Tokens"].(discord.TokenRequestResult)
	if !ok {
		return "", errors.New("tokens don't exist in this session")
	}

	return tokens.AccessToken, nil
}

func GetCurrentUser(ctx context.Context) (*discord.User, error) {
	session := GetSession(ctx)

	user, ok := session.Values["User"].(discord.User)
	if !ok {
		return nil, errors.New("the user hasn't been fetched for this session yet")
	}

	return &user, nil
}

func DestroySession(ctx context.Context) {
	session := GetSession(ctx)
	session.Options.MaxAge = -1
}
