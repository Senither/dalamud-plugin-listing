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

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Dalamud Plugin Listing")

	resp, err := client.Do(req)
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
			LastUpdatedAt: time.Now().Unix(),
		}

		state.UpsertRepository(repo)
	}
}
