package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/senither/dalamud-plugin-listing/http/renders"
)

var srv *http.Server

func SetupServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRequest)
	mux.HandleFunc("/favicon.ico", handleFavicon)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	srv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("Starting server on", "addr", addr)
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("Server has been shutdown gracefully")
				return
			}

			slog.Error("Server caught an unexpected error", "error", err)
			os.Exit(1)
		}
	}()
}

func ShutdownServer() {
	slog.Info("Shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server",
			"error", err,
		)
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/icons/favicon.ico")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.URL.Path) == "/webhook/github-release" && r.Method == "POST" {
		handleGitHubReleaseWebhook(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/plugins/") && r.Method == "GET" {
		handleRenderingPluginByName(w, r, r.URL.Path[9:])
		return
	} else if strings.HasPrefix(r.URL.Path, "/authors/") && r.Method == "GET" {
		handleRenderingPluginByAuthors(w, r, r.URL.Path[9:])
		return
	}

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

	ApplyHttpBaseResponseHeaders(w)

	if r.Header.Get("Accept") == "application/json" {
		renders.RenderJson(w, r)
		return
	}

	renders.RenderHtml(w, r)
}

func ApplyHttpBaseResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "accept,content-type")
	w.Header().Set("Access-Control-Max-Age", "300")

	w.Header().Set("Server", "'; DROP TABLE servertypes; --")
	w.Header().Set("X-Powered-By", "Nerd Rage and Caffeine")
	w.Header().Set("X-Accepts", "text/html,application/json")
}
