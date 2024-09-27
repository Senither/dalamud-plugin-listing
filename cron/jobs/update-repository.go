package jobs

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/senither/dalamud-plugin-listing/state"
)

type UpdateRepositoryJob struct {
	Interval     time.Duration
	RunOnStartup bool
	Ticker       *time.Ticker
}

var jobs = make(map[string]*UpdateRepositoryJob)

func StartUpdateRepositoryJob(url string, interval time.Duration, runOnStartup bool) {
	if runOnStartup {
		runUpdate(url)
	}

	tick := time.NewTicker(interval)

	jobs[url] = &UpdateRepositoryJob{
		Interval:     interval,
		RunOnStartup: runOnStartup,
		Ticker:       tick,
	}

	go func() {
		for range tick.C {
			runUpdate(url)
		}
	}()
}

func GetJobs() map[string]*UpdateRepositoryJob {
	return jobs
}

func runUpdate(url string) {
	slog.Info("Sending request to update repository for",
		"url", url,
	)

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

	repos, err := decodeJsonRequestBody(resp.Body)
	if err != nil {
		slog.Error("Failed to decode JSON response",
			"err", err,
			"url", url,
		)
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

func decodeJsonRequestBody(body io.ReadCloser) ([]state.Repository, error) {
	reqBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	reqBody := string(reqBytes)
	exp := regexp.MustCompile(`,(\s*[\}\]])`)
	reqBody = exp.ReplaceAllString(reqBody, "$1")

	var repos []state.Repository

	decoder := json.NewDecoder(bytes.NewBufferString(reqBody))
	err = decoder.Decode(&repos)

	if err != nil {
		return nil, err
	}
	return repos, nil
}
