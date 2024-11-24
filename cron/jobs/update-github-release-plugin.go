package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/senither/dalamud-plugin-listing/state"
)

type UpdatePluginReleaseJob struct {
	Interval     time.Duration
	RunOnStartup bool
	Ticker       *time.Ticker
}

type GitHubPluginReleaseResponse struct {
	Url        string `json:"url"`
	TagName    string `json:"tag_name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	Body       string `json:"body"`
	Assets     []GitHubPluginReleaseAsset
}

type GitHubPluginReleaseAsset struct {
	Url                string `json:"url"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

var jobs = make(map[string]*UpdatePluginReleaseJob)

func StartUpdatePluginReleaseJob(repoName string, interval time.Duration, runOnStartup bool) {
	if runOnStartup {
		runUpdatePluginRelease(repoName)
	}

	tick := time.NewTicker(interval)

	jobs[repoName] = &UpdatePluginReleaseJob{
		Interval:     interval,
		RunOnStartup: runOnStartup,
		Ticker:       tick,
	}

	go func() {
		for range tick.C {
			runUpdatePluginRelease(repoName)
		}
	}()
}

func RunGitHubReleaseUpdateJob(repoName string) {
	runUpdatePluginRelease(repoName)
}

func GetPluginReleasesJobs() map[string]*UpdatePluginReleaseJob {
	return jobs
}

func runUpdatePluginRelease(repoName string) {
	slog.Info("Sending request to update plugin release for",
		"repoName", repoName,
	)

	client := http.Client{}
	releaseReq, releaseErr := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repoName), nil)
	if releaseErr != nil {
		log.Fatal(releaseErr)
	}

	releaseReq.Header.Set("Accept", "application/json")
	releaseReq.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")

	releaseResp, releaseErr := client.Do(releaseReq)
	if releaseErr != nil {
		log.Fatal(releaseErr)
	}

	defer releaseResp.Body.Close()

	release, releaseErr := decodeJsonPluginReleaseRequestBody(releaseResp.Body)
	if releaseErr != nil {
		slog.Error("Failed to decode JSON response",
			"err", releaseErr,
			"repoName", repoName,
		)
		return
	}

	var manifestAsset, releaseAsset = getManifestAndLatestReleaseAssets(*release)
	if manifestAsset == nil || releaseAsset == nil {
		slog.Error("Failed to find a manifest or release asset in the release",
			"repoName", repoName,
			"release", releaseAsset,
			"manifest", manifestAsset,
		)
		return
	}

	assetReq, assetErr := http.NewRequest("GET", manifestAsset.BrowserDownloadUrl, nil)
	if assetErr != nil {
		log.Fatal(assetErr)
	}

	assetReq.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")

	manifestResp, assetErr := client.Do(assetReq)
	if assetErr != nil {
		log.Fatal(assetErr)
	}

	defer manifestResp.Body.Close()

	manifestBytes, assetErr := io.ReadAll(manifestResp.Body)
	if assetErr != nil {
		log.Fatal(assetErr)
	}

	var repository state.Repository
	manifestErr := json.Unmarshal(manifestBytes, &repository)
	if manifestErr != nil {
		slog.Error("Failed to decode JSON manifest",
			"err", manifestErr,
			"repoName", repoName,
		)
		return
	}

	var truthy = true

	var repoUrl = fmt.Sprintf("https://github.com/%s", repoName)
	var repositoryOrigin = state.RepositoryOrigin{
		LastUpdatedAt:    time.Now().Unix(),
		RepositoryUrl:    repoUrl,
		IsInternalPlugin: &truthy,
	}

	repository.RepoUrl = &repoUrl
	repository.DownloadLinkInstall = &releaseAsset.BrowserDownloadUrl
	repository.DownloadLinkUpdate = &releaseAsset.BrowserDownloadUrl
	repository.RepositoryOrigin = repositoryOrigin

	state.UpsertRepository(repository)
}

func getManifestAndLatestReleaseAssets(release GitHubPluginReleaseResponse) (*GitHubPluginReleaseAsset, *GitHubPluginReleaseAsset) {
	var manifestAsset *GitHubPluginReleaseAsset = nil
	var latestAsset *GitHubPluginReleaseAsset = nil

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, ".json") && asset.ContentType == "application/json" {
			manifestAsset = &asset
		} else if strings.Contains(asset.Name, ".zip") {
			latestAsset = &asset
		}
	}

	return manifestAsset, latestAsset
}

func decodeJsonPluginReleaseRequestBody(body io.ReadCloser) (*GitHubPluginReleaseResponse, error) {
	reqBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	reqBody := string(reqBytes)
	exp := regexp.MustCompile(`,(\s*[\}\]])`)
	reqBody = exp.ReplaceAllString(reqBody, "$1")

	var release GitHubPluginReleaseResponse

	decoder := json.NewDecoder(bytes.NewBufferString(reqBody))
	err = decoder.Decode(&release)

	if err != nil {
		return nil, err
	}
	return &release, nil
}
