package middlewares

import (
	"net/http"
	"strings"
)

func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") && !(r.URL.Path == "/") {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], 301)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
