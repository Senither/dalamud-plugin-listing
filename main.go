package main

import (
	"fmt"
	"net/http"

	"github.com/senither/dalamud-plugin-listing/renders"
)

func main() {
	http.HandleFunc("/", handleRequest)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renders.RenderError(w, r)
		return
	}

	w.Header().Set("Server", "'; DROP TABLE servertypes; --")
	w.Header().Set("X-Powered-By", "Nerd Rage and Caffeine")
	w.Header().Set("X-Accepts", "text/html,application/json")

	if r.Header.Get("Accept") == "application/json" {
		renders.RenderJson(w, r)
		return
	}

	renders.RenderHtml(w, r)
}
