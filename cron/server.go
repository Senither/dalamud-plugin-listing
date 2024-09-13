package cron

import (
	"math/rand/v2"
	"time"

	"github.com/senither/dalamud-plugin-listing/cron/jobs"
	"github.com/senither/dalamud-plugin-listing/state"
)

func SetupJobs() {
	state.LoadCachedRepositoryDataFromDisk()
	state.LoadRepositoriesFromDisk()

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

	jobs.StartDeleteExpiredRepositoriesJob(time.Second * 30)
}
