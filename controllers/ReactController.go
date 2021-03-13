package controllers

import (
	"net/http"
	"os"
	"path/filepath"
)

func ReactController(staticPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(filepath.Join(staticPath, r.URL.Path)); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.ServeFile(w, r, filepath.Join(staticPath, r.URL.Path))
		}
	}
}
