package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func ReactHandler(staticPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if a, err := os.Stat(filepath.Join(staticPath, r.URL.Path)); os.IsNotExist(err) {
			fmt.Println(a)
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

