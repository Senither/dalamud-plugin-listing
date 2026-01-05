package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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

var jobs = make(map[string]*UpdatePluginReleaseJob)

func StartUpdatePluginReleaseJob(repoName string, interval time.Duration, runOnStartup bool) {
	ip := state.GetInternalPluginByName(repoName)
	if ip == nil {
		slog.Error("Failed to find internal plugin for GitHub release update job",
			"repoName", repoName,
		)
		return
	}

	if runOnStartup {
		runUpdatePluginRelease(ip)
	}

	tick := time.NewTicker(interval)

	jobs[repoName] = &UpdatePluginReleaseJob{
		Interval:     interval,
		RunOnStartup: runOnStartup,
		Ticker:       tick,
	}

	go func() {
		for range tick.C {
			runUpdatePluginRelease(ip)
		}
	}()
}

func RunGitHubReleaseUpdateJob(repoName string) {
	ip := state.GetInternalPluginByName(repoName)
	if ip == nil {
		slog.Error("Failed to find internal plugin for GitHub release update job",
			"repoName", repoName,
		)
		return
	}

	runUpdatePluginRelease(ip)
}

func GetPluginReleasesJobs() map[string]*UpdatePluginReleaseJob {
	return jobs
}

func runUpdatePluginRelease(ip *state.InternalPlugin) {
	slog.Info("Sending request to update plugin release for",
		"repoName", ip.Name,
		"private", ip.Private,
	)

	client := http.Client{}
	releaseReq, releasesErr := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=100", ip.Name), nil)
	if releasesErr != nil {
		slog.Error("Failed to create plugin release request",
			"err", releasesErr,
			"repoName", ip.Name,
		)
		return
	}

	releaseReq.Header.Set("Accept", "application/json")
	releaseReq.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" && ip.Private {
		releaseReq.Header.Set("Authorization", "Bearer "+token)
	}

	releaseResp, releasesErr := client.Do(releaseReq)
	if releasesErr != nil {
		slog.Error("Failed to communicate with GitHub API",
			"err", releasesErr,
			"repoName", ip.Name,
		)
		return
	}

	defer releaseResp.Body.Close()

	releases, releasesErr := decodeJsonPluginReleaseRequestBody(releaseResp.Body)
	if releasesErr != nil {
		slog.Error("Failed to decode JSON response",
			"err", releasesErr,
			"repoName", ip.Name,
		)
		return
	}

	if len(releases) == 0 {
		slog.Error("Failed to find any releases for repository",
			"repoName", ip.Name,
		)
		return
	}

	state.UpsertReleaseMetadata(ip.Name, releases)

	var manifestAsset, releaseAsset = state.GetManifestAndLatestReleaseAssets(releases[0])
	if manifestAsset == nil || releaseAsset == nil {
		slog.Error("Failed to find a manifest or release asset in the release",
			"repoName", ip.Name,
			"release", releaseAsset,
			"manifest", manifestAsset,
		)
		return
	}

	manifestUrl := manifestAsset.BrowserDownloadUrl
	if ip.Private {
		manifestUrl = manifestAsset.Url
	}

	assetReq, assetErr := http.NewRequest("GET", manifestUrl, nil)
	if assetErr != nil {
		slog.Error("Failed to create asset request",
			"err", assetErr,
			"repoName", ip.Name,
			"manifestAsset", manifestAsset,
			"downloadUrl", manifestAsset.BrowserDownloadUrl,
		)
		return
	}

	assetReq.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" && ip.Private {
		assetReq.Header.Set("Authorization", "Bearer "+token)
		assetReq.Header.Set("Accept", "application/octet-stream")
	}

	manifestResp, assetErr := client.Do(assetReq)
	if assetErr != nil {
		slog.Error("Failed to communicate with asset URL",
			"err", assetErr,
			"repoName", ip.Name,
			"downloadUrl", manifestAsset.BrowserDownloadUrl,
		)
		return
	}

	defer manifestResp.Body.Close()

	manifestBytes, assetErr := io.ReadAll(manifestResp.Body)
	if assetErr != nil {
		slog.Error("Failed to read asset response body",
			"err", assetErr,
			"repoName", ip.Name,
			"downloadUrl", manifestAsset.BrowserDownloadUrl,
		)
		return
	}

	var repository state.Repository
	manifestErr := json.Unmarshal(manifestBytes, &repository)
	if manifestErr != nil {
		slog.Error("Failed to decode JSON manifest",
			"err", manifestErr,
			"repoName", ip.Name,
		)
		return
	}

	var truthy = true

	var repoUrl = fmt.Sprintf("https://github.com/%s", ip.Name)
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

func decodeJsonPluginReleaseRequestBody(body io.ReadCloser) ([]state.GitHubPluginRelease, error) {
	reqBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	reqBody := string(reqBytes)
	exp := regexp.MustCompile(`,(\s*[\}\]])`)
	reqBody = exp.ReplaceAllString(reqBody, "$1")

	var release []state.GitHubPluginRelease

	decoder := json.NewDecoder(bytes.NewBufferString(reqBody))
	err = decoder.Decode(&release)

	if err != nil {
		return nil, err
	}

	return release, nil
}
