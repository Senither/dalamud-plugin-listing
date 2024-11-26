package http

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/senither/dalamud-plugin-listing/state"
)

type RepositorySearchCallback func(repo *state.Repository, searchQuery string) bool

func handleRenderingPluginByInternalName(w http.ResponseWriter, r *http.Request, internalName string) {
	renderPluginRepositorySearch(w, r, internalName, func(repo *state.Repository, searchQuery string) bool {
		return strings.HasPrefix(strings.ToLower(repo.InternalName), searchQuery)
	})
}

func handleRenderingPluginByAuthors(w http.ResponseWriter, r *http.Request, author string) {
	renderPluginRepositorySearch(w, r, author, func(repo *state.Repository, searchQuery string) bool {
		return strings.HasPrefix(strings.ToLower(repo.Author), searchQuery)
	})
}

func renderPluginRepositorySearch(w http.ResponseWriter, r *http.Request, value string, callback RepositorySearchCallback) {
	ApplyHttpBaseResponseHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	searchQuery := strings.ToLower(value)
	searchQuery = strings.Trim(searchQuery, "/")
	searchQuery = strings.TrimSuffix(searchQuery, ".json")

	slog.Info("Handling search request",
		"searchQuery", searchQuery,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	var repositories []*state.Repository
	for _, repo := range state.GetRepositories() {
		if callback(&repo, searchQuery) {
			repositories = append(repositories, &repo)
		}
	}

	if len(repositories) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status": 404, "reason": "No plugin(s) matching with the given search parameter were found"}`))
		return
	}

	content, err := json.Marshal(repositories)
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
