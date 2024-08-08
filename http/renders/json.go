package renders

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/senither/dalamud-plugin-listing/state"
)

func RenderJson(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" - Rendering JSON")

	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(state.GetRepositories())
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	fmt.Fprintf(w, string(content))
}
