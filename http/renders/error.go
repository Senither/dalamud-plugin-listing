package renders

import (
	"net/http"

	"github.com/senither/dalamud-plugin-listing/metrics"
)

func RenderError(w http.ResponseWriter, r *http.Request) {
	metrics.IncrementRouteRequestCounter(metrics.ErrorRoute)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
