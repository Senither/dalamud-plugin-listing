package renders

import (
	"net/http"
)

func RenderError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
