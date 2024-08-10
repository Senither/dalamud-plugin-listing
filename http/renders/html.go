package renders

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/senither/dalamud-plugin-listing/metrics"
)

func RenderHtml(w http.ResponseWriter, r *http.Request) {
	metrics.IncrementRouteRequestCounter(metrics.HtmlRoute)

	fmt.Println(" - Rendering HTML")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	content, err := os.ReadFile("./views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(content))
}
