package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/senither/dalamud-plugin-listing/state"
)

func UpdateRepository(url string, interval time.Duration, runOnStartup bool) {
	if runOnStartup {
		run(url)
	}

	tick := time.NewTicker(interval)

	go func() {
		for range tick.C {
			run(url)
		}
	}()
}

func run(url string) {
	fmt.Println("Sending request to update repository for: ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var repos []state.Repository

	// Parse the JSON array from the request body
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&repos)
	if err != nil {
		fmt.Println("Failed to decode JSON response:", err)
		return
	}

	for _, repo := range repos {
		repo.RepositoryOrigin = state.RepositoryOrigin{
			RepositoryUrl: url,
			LastUpdatedAt: time.Now().Format(time.RFC3339),
		}

		state.UpsertRepository(repo)
	}
}
