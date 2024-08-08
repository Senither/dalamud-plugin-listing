package cron

import (
	"time"

	"github.com/senither/dalamud-plugin-listing/cron/jobs"
	"github.com/senither/dalamud-plugin-listing/state"
)

func SetupJobs() {
	// Loops through all the repositories in the state and creates a new job for each one.
	for _, repo := range state.GetRepositories() {
		jobs.UpdateRepository(repo, time.Second*5, true)
	}
}
