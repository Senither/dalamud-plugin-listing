package state

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type Repository struct {
	Author              string           `json:"Author"`
	Name                string           `json:"Name"`
	Punchline           *string          `json:"Punchline,omitempty"`
	Description         string           `json:"Description"`
	Changelog           *string          `json:"Changelog,omitempty"`
	InternalName        string           `json:"InternalName"`
	AssemblyVersion     interface{}      `json:"AssemblyVersion,omitempty"`
	RepoUrl             string           `json:"RepoUrl"`
	IconUrl             string           `json:"IconUrl"`
	ApplicableVersion   *string          `json:"ApplicableVersion,omitempty"`
	Tags                []string         `json:"Tags"`
	DalamudApiLevel     interface{}      `json:"DalamudApiLevel,omitempty"`
	IsHide              *interface{}     `json:"IsHide,omitempty"`
	IsTestingExclusive  *interface{}     `json:"IsTestingExclusive,omitempty"`
	LastUpdated         interface{}      `json:"LastUpdated,omitempty"`
	DownloadCount       interface{}      `json:"DownloadCount,omitempty"`
	DownloadLinkInstall string           `json:"DownloadLinkInstall"`
	DownloadLinkTesting *string          `json:"DownloadLinkTesting,omitempty"`
	DownloadLinkUpdate  *string          `json:"DownloadLinkUpdate,omitempty"`
	RepositoryOrigin    RepositoryOrigin `json:"OriginRepositoryUrl"`
}

type RepositoryOrigin struct {
	RepositoryUrl string `json:"RepositoryUrl"`
	LastUpdatedAt int64  `json:"LastUpdatedAt"`
}

var repositories []Repository
var timer = time.NewTimer(time.Nanosecond)

func UpsertRepository(repo Repository) {
	index := getRepositoryIndex(repo.RepoUrl)

	if index == -1 {
		repositories = append(repositories, repo)
	} else {
		repositories[index] = repo
	}

	writeRepositoriesToDisk()
}

func DeleteRepository(url string) {
	index := getRepositoryIndex(url)
	if index == -1 {
		return
	}

	repositories = append(repositories[:index], repositories[index+1:]...)

	writeRepositoriesToDisk()
}

func GetRepositories() []Repository {
	return repositories
}

func GetRepositoriesByOriginUrl(url string) []Repository {
	var filteredRepos []Repository

	for _, repository := range repositories {
		if repository.RepositoryOrigin.RepositoryUrl == url {
			filteredRepos = append(filteredRepos, repository)
		}
	}

	return filteredRepos
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

		AddUrl(repo)
	}
}

func getRepositoryIndex(url string) int {
	for i, repository := range repositories {
		if repository.RepoUrl == url {
			return i
		}
	}

	return -1
}

func writeRepositoriesToDisk() {
	if timer != nil {
		timer.Stop()
	}

	timer = time.AfterFunc(5*time.Second, func() {
		content, err := json.Marshal(GetRepositories())
		if err != nil {
			log.Fatalf("Error converting to JSON: %v", err)
		}

		os.WriteFile("./cached-repositories.json", content, 0644)
	})
}
