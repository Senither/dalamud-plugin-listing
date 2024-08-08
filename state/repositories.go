package state

type Repository struct {
	Author              string       `json:"Author"`
	Name                string       `json:"Name"`
	Punchline           *string      `json:"Punchline"`
	Description         string       `json:"Description"`
	Changelog           *string      `json:"Changelog"`
	InternalName        string       `json:"InternalName"`
	AssemblyVersion     interface{}  `json:"AssemblyVersion"`
	RepoUrl             string       `json:"RepoUrl"`
	IconUrl             string       `json:"IconUrl"`
	ApplicableVersion   *string      `json:"ApplicableVersion"`
	Tags                []string     `json:"Tags"`
	DalamudApiLevel     interface{}  `json:"DalamudApiLevel"`
	IsHide              *interface{} `json:"IsHide"`
	IsTestingExclusive  *interface{} `json:"IsTestingExclusive"`
	LastUpdated         interface{}  `json:"LastUpdated"`
	DownloadCount       interface{}  `json:"DownloadCount"`
	DownloadLinkInstall string       `json:"DownloadLinkInstall"`
	DownloadLinkTesting *string      `json:"DownloadLinkTesting"`
	DownloadLinkUpdate  *string      `json:"DownloadLinkUpdate"`
}

var repositories []Repository

func UpsertRepository(repo Repository) {
	index := getRepositoryIndex(repo.RepoUrl)

	if index == -1 {
		repositories = append(repositories, repo)
	} else {
		repositories[index] = repo
	}
}

func GetRepositories() []Repository {
	return repositories
}

func getRepositoryIndex(url string) int {
	for i, repository := range repositories {
		if repository.RepoUrl == url {
			return i
		}
	}

	return -1
}
