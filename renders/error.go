package renders

import (
	"fmt"
	"net/http"
)

func RenderError(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" - 404 Not Found")

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
