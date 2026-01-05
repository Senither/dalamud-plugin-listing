package state

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type Repository struct {
	Author                 string           `json:"Author"`
	Name                   string           `json:"Name"`
	Punchline              *string          `json:"Punchline,omitempty"`
	Description            string           `json:"Description"`
	Changelog              *string          `json:"Changelog,omitempty"`
	InternalName           string           `json:"InternalName"`
	AssemblyVersion        interface{}      `json:"AssemblyVersion,omitempty"`
	TestingAssemblyVersion interface{}      `json:"TestingAssemblyVersion,omitempty"`
	RepoUrl                *string          `json:"RepoUrl"`
	IconUrl                *string          `json:"IconUrl,omitempty"`
	ApplicableVersion      *string          `json:"ApplicableVersion,omitempty"`
	Tags                   []string         `json:"Tags"`
	DalamudApiLevel        interface{}      `json:"DalamudApiLevel,omitempty"`
	IsOutdated             bool             `json:"IsOutdated"`
	TestingDalamudApiLevel interface{}      `json:"TestingDalamudApiLevel,omitempty"`
	IsHide                 *interface{}     `json:"IsHide,omitempty"`
	IsTestingExclusive     *interface{}     `json:"IsTestingExclusive,omitempty"`
	LastUpdated            interface{}      `json:"LastUpdated,omitempty"`
	DownloadCount          interface{}      `json:"DownloadCount,omitempty"`
	LastUpdate             interface{}      `json:"LastUpdate,omitempty"`
	LoadPriority           *int64           `json:"LoadPriority,omitempty"`
	LoadRequiredState      *int64           `json:"LoadRequiredState,omitempty"`
	LoadSync               *bool            `json:"LoadSync,omitempty"`
	AcceptsFeedback        *bool            `json:"AcceptsFeedback,omitempty"`
	DownloadLinkInstall    *string          `json:"DownloadLinkInstall,omitempty"`
	DownloadLinkTesting    *string          `json:"DownloadLinkTesting,omitempty"`
	DownloadLinkUpdate     *string          `json:"DownloadLinkUpdate,omitempty"`
	RepositoryOrigin       RepositoryOrigin `json:"OriginRepositoryUrl"`
}

type RepositoryOrigin struct {
	RepositoryUrl    string `json:"RepositoryUrl"`
	LastUpdatedAt    int64  `json:"LastUpdatedAt"`
	IsInternalPlugin *bool  `json:"IsInternalPlugin,omitempty"`
	IsPrivatePlugin  *bool  `json:"IsPrivatePlugin,omitempty"`
}

var (
	repositories            []Repository
	repositoryTimer         = time.NewTimer(time.Nanosecond)
	repositoryLastUpdatedAt = time.Now().Unix()
)

func UpsertRepository(repo Repository) {
	if repo.RepoUrl == nil || *repo.RepoUrl == "" {
		repo.RepoUrl = findRepositoryUrl(repo)
	}

	index := getRepositoryIndex(repo)

	if index == -1 {
		repositories = append(repositories, repo)
	} else {
		repositories[index] = repo
	}

	writeRepositoriesToDisk()
}

func DeleteRepository(repo Repository) {
	index := getRepositoryIndex(repo)
	if index == -1 {
		return
	}

	repositories = append(repositories[:index], repositories[index+1:]...)

	writeRepositoriesToDisk()
}

func GetRepositories() []Repository {
	latestDalamudApiLevel := GetLatestDalamudApiLevel()

	for i, repository := range repositories {
		if repository.DalamudApiLevel != nil {
			level, ok := repository.DalamudApiLevel.(float64)
			if ok {
				repositories[i].IsOutdated = level != latestDalamudApiLevel
			}
		}
	}

	return repositories
}

func GetRepositoriesSize() int {
	return len(repositories)
}

func GetRepositoriesByOriginUrl(url string) []Repository {
	var filteredRepos []Repository

	for _, repository := range GetRepositories() {
		if repository.RepositoryOrigin.RepositoryUrl == url {
			filteredRepos = append(filteredRepos, repository)
		}
	}

	return filteredRepos
}

func GetRepositoryByGitHubReleaseRepositoryName(repoName string) *Repository {
	githubLink := "github.com/" + repoName
	localLink := fmt.Sprintf(
		"%s/download/%s",
		strings.TrimSuffix(strings.TrimSpace(os.Getenv("APP_URL")), "/"),
		repoName,
	)

	for _, repository := range GetRepositories() {
		ip := GetInternalPluginByName(repoName)
		if ip == nil {
			continue
		}

		downloadLink := getAvailableDownloadLink(repository)

		if ip.Private && strings.Contains(*downloadLink, localLink) {
			return &repository
		} else if !ip.Private && strings.Contains(*downloadLink, githubLink) {
			return &repository
		}
	}

	return nil
}

func GetRepositoriesLastUpdatedAt() int64 {
	return repositoryLastUpdatedAt
}

func LoadCachedRepositoryDataFromDisk() {
	content, err := os.ReadFile("cached-repositories.json")
	if err != nil {
		return
	}

	var repositories []Repository
	if err := json.Unmarshal(content, &repositories); err != nil {
		log.Fatalf("Error converting JSON: %v", err)
	}

	for _, repo := range repositories {
		UpsertRepository(repo)
	}
}

func LoadRepositoriesFromDisk() {
	content, err := os.ReadFile("repositories.txt")
	if err != nil {
		log.Fatal(err)
	}

	repositories := strings.Split(string(content), "\n")

	for _, repo := range repositories {
		if repo == "" {
			continue
		}

		AddUrl(strings.Trim(repo, "\r"))
	}
}

func LoadPluginsFromDisk() {
	content, err := os.ReadFile("plugins.txt")
	if err != nil {
		log.Fatal(err)
	}

	plugins := strings.Split(string(content), "\n")

	for _, repo := range plugins {
		if repo == "" {
			continue
		}

		AddInternalPluginUrl(strings.Trim(repo, "\r"))
	}
}

func GetLatestDalamudApiLevel() float64 {
	var latestDalamudApiLevel float64

	for _, repo := range repositories {
		if repo.DalamudApiLevel == nil {
			continue
		}

		level, ok := repo.DalamudApiLevel.(float64)
		if ok && level > latestDalamudApiLevel {
			latestDalamudApiLevel = level
		}
	}

	return latestDalamudApiLevel
}

func getRepositoryIndex(repo Repository) int {
	for i, repository := range repositories {
		if repository.Name == repo.Name &&
			repository.Author == repo.Author &&
			repository.InternalName == repo.InternalName {
			return i
		}
	}

	return -1
}

func findRepositoryUrl(repo Repository) *string {
	downloadLink := getAvailableDownloadLink(repo)
	if downloadLink == nil {
		return nil
	}

	parsedUrl, err := url.ParseRequestURI(*downloadLink)
	if err != nil {
		return nil
	}

	if !strings.HasPrefix(*downloadLink, "https://github.com") {
		url := parsedUrl.Scheme + "://" + parsedUrl.Host
		return &url
	}

	pathParts := strings.Split(parsedUrl.Path, "/")
	if len(pathParts) < 3 {
		return nil
	}

	url := parsedUrl.Scheme + "://" + parsedUrl.Host + "/" + pathParts[1] + "/" + pathParts[2]
	return &url
}

func getAvailableDownloadLink(repo Repository) *string {
	if repo.DownloadLinkInstall != nil && *repo.DownloadLinkInstall != "" {
		return repo.DownloadLinkInstall
	}

	if repo.DownloadLinkTesting != nil && *repo.DownloadLinkTesting != "" {
		return repo.DownloadLinkTesting
	}

	if repo.DownloadLinkUpdate != nil && *repo.DownloadLinkUpdate != "" {
		return repo.DownloadLinkUpdate
	}

	return nil
}

func writeRepositoriesToDisk() {
	repositoryLastUpdatedAt = time.Now().Unix()

	if repositoryTimer != nil {
		repositoryTimer.Stop()
	}

	repositoryTimer = time.AfterFunc(5*time.Second, func() {
		content, err := json.Marshal(repositories)
		if err != nil {
			log.Fatalf("Error converting to JSON: %v", err)
		}

		os.WriteFile("./cached-repositories.json", content, 0644)
	})
}
