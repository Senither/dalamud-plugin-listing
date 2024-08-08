package state

var repositories []string

func AddRepository(url string) {
	if exists(url) {
		return
	}

	repositories = append(repositories, url)
}

func GetRepositories() []string {
	return repositories
}

func exists(url string) bool {
	for _, repo := range repositories {
		if repo == url {
			return true
		}
	}

	return false
}
