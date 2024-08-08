package renders

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func RenderJson(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" - Rendering JSON")

	w.Header().Set("Content-Type", "application/json")

	content, err := os.ReadFile("./views/api-response.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(content))
}
