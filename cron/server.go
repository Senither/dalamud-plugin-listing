package cron

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/senither/dalamud-plugin-listing/cron/jobs"
	"github.com/senither/dalamud-plugin-listing/state"
)

func SetupJobs() {
	loadRepositories()

	// Loops through all the repositories in the state and creates a new job for each one.
	for _, repo := range state.GetUrls() {
		jobs.UpdateRepository(repo, time.Minute*30, true)
	}
}

func loadRepositories() {
	content, err := os.ReadFile("repositories.txt")
	if err != nil {
		log.Fatal(err)
	}

	repositories := strings.Split(string(content), "\n")

	for _, repo := range repositories {
		if repo == "" {
			continue
		}

		state.AddUrl(repo)
	}
}
