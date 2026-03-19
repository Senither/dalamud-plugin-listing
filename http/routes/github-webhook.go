package routes

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
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

func GitHubReleaseWebhook(c fiber.Ctx) error {
	slog.InfoContext(c, "Handling GitHub release webhook",
		"remote", c.IP(),
	)

	var req GitHubWebhookRequest = GitHubWebhookRequest{}
	if err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&req); err != nil {
		slog.Error("Failed to decode GitHub release webhook request",
			"error", err,
		)

		return c.Status(http.StatusBadRequest).SendString("Failed to decode request")
	}

	for _, internalPlugin := range state.GetInternalPlugins() {
		if internalPlugin.Name == req.Repository.FullName {
			go func() {
				slog.Info("Running GitHub release update job in 10 seconds",
					"repository", req.Repository.FullName,
				)

				// Run the job in 10 seconds to give GitHub time to process the release
				// and make it available for the job to fetch.
				<-time.After(10 * time.Second)

				jobs.RunGitHubReleaseUpdateJob(req.Repository.FullName)
			}()

			return c.SendStatus(fiber.StatusAccepted)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}
