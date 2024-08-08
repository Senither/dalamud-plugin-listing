package http

import (
	"fmt"
	"net/http"

	"github.com/senither/dalamud-plugin-listing/http/renders"
)

func SetupServer() {
	http.HandleFunc("/", handleRequest)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Method + ":" + r.URL.Path + " from " + r.RemoteAddr)

	if r.URL.Path != "/" {
		fmt.Println("404 Not Found:", r.URL.Path)

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
