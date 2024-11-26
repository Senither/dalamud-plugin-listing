package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/senither/dalamud-plugin-listing/state"
)

type RepositorySearchCallback func(repo *state.Repository, searchQuery string) bool

func handleRenderingPluginByInternalName(w http.ResponseWriter, internalName string) {
	renderPluginRepositorySearch(w, internalName, func(repo *state.Repository, searchQuery string) bool {
		return strings.HasPrefix(strings.ToLower(repo.InternalName), searchQuery)
	})
}

func handleRenderingPluginByAuthors(w http.ResponseWriter, author string) {
	renderPluginRepositorySearch(w, author, func(repo *state.Repository, searchQuery string) bool {
		return strings.HasPrefix(strings.ToLower(repo.Author), searchQuery)
	})
}

func renderPluginRepositorySearch(w http.ResponseWriter, value string, callback RepositorySearchCallback) {
	ApplyHttpBaseResponseHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	searchQuery := strings.ToLower(value)
	searchQuery = strings.Trim(searchQuery, "/")
	searchQuery = strings.TrimSuffix(searchQuery, ".json")

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
