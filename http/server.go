package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/senither/dalamud-plugin-listing/http/renders"
)

func SetupServer() {
	http.HandleFunc("/", handleRequest)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	slog.Info("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		slog.InfoContext(r.Context(), "Received request to invalid route",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)

		renders.RenderError(w, r)
		return
	}

	requestType := "html"
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		requestType = "json"
	}

	slog.InfoContext(r.Context(), "Handling request",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
		"requestType", requestType,
	)

	w.Header().Set("Server", "'; DROP TABLE servertypes; --")
	w.Header().Set("X-Powered-By", "Nerd Rage and Caffeine")
	w.Header().Set("X-Accepts", "text/html,application/json")

	if r.Header.Get("Accept") == "application/json" {
		renders.RenderJson(w, r)
		return
	}

	renders.RenderHtml(w, r)
}
