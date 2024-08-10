package renders

import (
	"fmt"
	"net/http"

	"github.com/senither/dalamud-plugin-listing/metrics"
)

func RenderError(w http.ResponseWriter, r *http.Request) {
	metrics.IncrementRouteRequestCounter(metrics.ErrorRoute)

	fmt.Println(" - 404 Not Found")

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
