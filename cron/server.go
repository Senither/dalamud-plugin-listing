package cron

import (
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
			if repo.RepositoryOrigin.LastUpdatedAt > time.Now().Add(time.Minute*15*-1).Unix() {
				runOnStart = false
			}
		}

		jobs.StartUpdateRepositoryJob(repoUrl, time.Minute*30, runOnStart)
	}

	jobs.StartDeleteExpiredRepositoriesJob(time.Second * 30)
}
