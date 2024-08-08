package main

import (
	"log"
	"os"
	"strings"

	"github.com/senither/dalamud-plugin-listing/cron"
	"github.com/senither/dalamud-plugin-listing/http"
	"github.com/senither/dalamud-plugin-listing/state"
)

func main() {
	loadRepositories()
	runApplication()
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

		state.AddRepository(repo)
	}
}

func runApplication() {
	cron.SetupJobs()
	http.SetupServer()
}
