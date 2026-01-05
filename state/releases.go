package state

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type GitHubReleaseContext struct {
	RepositoryName string
	Releases       []GitHubPluginRelease
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

var (
	releaseContexts       []GitHubReleaseContext
	releasesTimer         = time.NewTimer(time.Nanosecond)
	releasesLastUpdatedAt = time.Now().Unix()
)

func UpsertReleaseMetadata(repoName string, releases []GitHubPluginRelease) {
	ip := GetInternalPluginByName(repoName)
	if ip == nil {
		log.Printf("Failed to find internal plugin for release metadata upsert: %s", repoName)
		return
	}

	var index = -1

	for i, r := range releaseContexts {
		if r.RepositoryName == repoName {
			index = i
			break
		}
	}

	if index == -1 {
		releaseContexts = append(releaseContexts, GitHubReleaseContext{
			RepositoryName: ip.Name,
			Releases:       releases,
		})
	} else {
		releaseContexts[index] = GitHubReleaseContext{
			RepositoryName: ip.Name,
			Releases:       releases,
		}
	}

	writePluginReleasesToDisk()
}

func GetDownloadUrlForPrivatePlugin(repoName string, tag string, asset *GitHubPluginReleaseAsset) string {
	url := strings.TrimSuffix(strings.TrimSpace(os.Getenv("APP_URL")), "/")

	return fmt.Sprintf("%s/download/%s/%s/%s", url, repoName, tag, asset.Name)
}

func GetReleaseMetadataByRepositoryName(repoName string) *GitHubReleaseContext {
	for _, r := range releaseContexts {
		if r.RepositoryName == repoName {
			return &r
		}
	}

	return nil
}

func GetManifestAndLatestReleaseAssets(release GitHubPluginRelease) (*GitHubPluginReleaseAsset, *GitHubPluginReleaseAsset) {
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

func LoadCachedPluginReleasesDataFromDisk() {
	content, err := os.ReadFile("cached-plugin-releases.json")
	if err != nil {
		return
	}

	var repositories []GitHubReleaseContext
	if err := json.Unmarshal(content, &repositories); err != nil {
		log.Fatalf("Error converting JSON: %v", err)
	}

	for _, repo := range repositories {
		UpsertReleaseMetadata(repo.RepositoryName, repo.Releases)
	}
}

func writePluginReleasesToDisk() {
	releasesLastUpdatedAt = time.Now().Unix()

	if releasesTimer != nil {
		releasesTimer.Stop()
	}

	releasesTimer = time.AfterFunc(5*time.Second, func() {
		content, err := json.Marshal(releaseContexts)
		if err != nil {
			log.Fatalf("Error converting to JSON: %v", err)
		}

		os.WriteFile("./cached-plugin-releases.json", content, 0644)
	})
}
