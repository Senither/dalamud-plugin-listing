package state

type Repository struct {
	Author              string
	Name                string
	Punchline           *string
	Description         string
	Changelog           *string
	InternalName        string
	AssemblyVersion     string
	RepoUrl             string
	IconUrl             string
	ApplicableVersion   *string
	Tags                []string
	DalamudApiLevel     int
	IsHide              *string
	IsTestingExclusive  *string
	LastUpdated         uint64
	DownloadCount       uint64
	DownloadLinkInstall string
	DownloadLinkTesting *string
	DownloadLinkUpdate  *string
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
