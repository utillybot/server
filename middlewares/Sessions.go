package middlewares

import (
	"context"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	"github.com/mitchellh/mapstructure"
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

func GetSessionData(ctx context.Context) (*helpers.SessionData, error) {
	session := ctx.Value(contextKeySession).(*sessions.Session)
	result := helpers.SessionData{}
	err := mapstructure.Decode(session.Values, &result)
	return &result, err
}

func GetSession(ctx context.Context) *sessions.Session {
	session := ctx.Value(contextKeySession).(*sessions.Session)
	return session
}

func DestroySession(ctx context.Context) {
	session := ctx.Value(contextKeySession).(*sessions.Session)
	session.Options.MaxAge = -1
}