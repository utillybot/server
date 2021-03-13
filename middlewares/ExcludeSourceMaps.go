package middlewares

import (
	"net/http"
	"strings"
)

func ExcludeSourceMaps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUser(r.Context())

		if (err == nil && user != nil && user.ID == "236279900728721409") || !strings.HasSuffix(r.URL.Path, ".map") {
			next.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
