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

type GitHubPluginRelease struct {
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
	DownloadCount      int    `json:"download_count"`
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
	releaseReq, releasesErr := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=100", repoName), nil)
	if releasesErr != nil {
		log.Fatal(releasesErr)
	}

	releaseReq.Header.Set("Accept", "application/json")
	releaseReq.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")

	releaseResp, releasesErr := client.Do(releaseReq)
	if releasesErr != nil {
		log.Fatal(releasesErr)
	}

	defer releaseResp.Body.Close()

	releases, releasesErr := decodeJsonPluginReleaseRequestBody(releaseResp.Body)
	if releasesErr != nil {
		slog.Error("Failed to decode JSON response",
			"err", releasesErr,
			"repoName", repoName,
		)
		return
	}

	if len(releases) == 0 {
		slog.Error("Failed to find any releases for repository",
			"repoName", repoName,
		)
		return
	}

	var manifestAsset, releaseAsset = getManifestAndLatestReleaseAssets(releases[0])
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

	totalDownloadCount := 0
	for _, release := range releases {
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, ".zip") {
				totalDownloadCount += asset.DownloadCount
			}
		}
	}

	repository.RepoUrl = &repoUrl
	repository.DownloadLinkInstall = &releaseAsset.BrowserDownloadUrl
	repository.DownloadLinkUpdate = &releaseAsset.BrowserDownloadUrl
	repository.RepositoryOrigin = repositoryOrigin
	repository.DownloadCount = totalDownloadCount

	state.UpsertRepository(repository)
}

func getManifestAndLatestReleaseAssets(release GitHubPluginRelease) (*GitHubPluginReleaseAsset, *GitHubPluginReleaseAsset) {
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

func decodeJsonPluginReleaseRequestBody(body io.ReadCloser) ([]GitHubPluginRelease, error) {
	reqBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	reqBody := string(reqBytes)
	exp := regexp.MustCompile(`,(\s*[\}\]])`)
	reqBody = exp.ReplaceAllString(reqBody, "$1")

	var release []GitHubPluginRelease

	decoder := json.NewDecoder(bytes.NewBufferString(reqBody))
	err = decoder.Decode(&release)

	if err != nil {
		return nil, err
	}

	return release, nil
}
