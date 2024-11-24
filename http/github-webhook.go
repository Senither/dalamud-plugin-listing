package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/senither/dalamud-plugin-listing/cron/jobs"
	"github.com/senither/dalamud-plugin-listing/state"
)

type GitHubWebhookRequest struct {
	HookId     int64                   `json:"hook_id"`
	Repository GitHubWebhookRepository `json:"repository"`
}

type GitHubWebhookRepository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func handleGitHubReleaseWebhook(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Handling GitHub release webhook",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)

	var req GitHubWebhookRequest = GitHubWebhookRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode GitHub release webhook request",
			"error", err,
		)
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	for _, repoName := range state.GetInternalPlugins() {
		if repoName == req.Repository.FullName {
			go func() {
				slog.Info("Running GitHub release update job in 10 seconds",
					"repository", req.Repository.FullName,
				)

				// Run the job in 10 seconds to give GitHub time to process the release
				// and make it available for the job to fetch.
				<-time.After(10 * time.Second)

				jobs.RunGitHubReleaseUpdateJob(req.Repository.FullName)
			}()

			w.WriteHeader(http.StatusAccepted)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
