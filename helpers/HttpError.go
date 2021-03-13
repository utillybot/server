package helpers

import "net/http"

func HttpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
