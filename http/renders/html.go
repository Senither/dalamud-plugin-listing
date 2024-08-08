package renders

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func RenderHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" - Rendering HTML")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	content, err := os.ReadFile("./views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(content))
}
