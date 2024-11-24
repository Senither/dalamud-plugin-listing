package cron

import (
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/senither/dalamud-plugin-listing/cron/jobs"
	"github.com/senither/dalamud-plugin-listing/state"
)

func SetupJobs() {
	state.LoadCachedRepositoryDataFromDisk()
	state.LoadRepositoriesFromDisk()
	state.LoadPluginsFromDisk()

	// Loops through all the repositories in the state and creates a new job for each one.
	for _, repoUrl := range state.GetUrls() {
		repos := state.GetRepositoriesByOriginUrl(repoUrl)

		runOnStart := true

		for _, repo := range repos {
			if repo.RepositoryOrigin.LastUpdatedAt > time.Now().Add(time.Minute*35*-1).Unix() {
				runOnStart = false
			}
		}

		// Gets a random number between 55 and 70 to add some randomness to
		// the job delay, then starts the job with the specified values.
		jobDelay := rand.IntN(15) + 55
		jobs.StartUpdateRepositoryJob(repoUrl, time.Minute*time.Duration(jobDelay), runOnStart)
	}

	for _, repoName := range state.GetInternalPlugins() {
		repo := state.GetRepositoryByGitHubReleaseRepositoryName(repoName)
		if repo == nil {
			slog.Warn("Failed to find plugin", "repo", repoName)
			// TODO: Start the job right away if the repo isn't found
			continue
		}

		slog.Info("Found plugin", "repo", repoName, "url", repo.RepositoryOrigin.RepositoryUrl)

		runOnStart := repo.RepositoryOrigin.LastUpdatedAt <= time.Now().Add(time.Minute*35*-1).Unix()

		// TODO: Starts the job and either runs it right away, or with a delay around 3 - 6 hours.
		slog.Info("Starting job",
			"repo", repoName,
			"url", repo.RepositoryOrigin.RepositoryUrl,
			"runOnStart", runOnStart,
		)
	}

	jobs.StartDeleteExpiredRepositoriesJob(time.Second * 30)
}

func ShutdownJobs() {
	for url, job := range jobs.GetJobs() {
		slog.Debug("Shutting down job", "url", url)

		job.Ticker.Stop()
	}
}
