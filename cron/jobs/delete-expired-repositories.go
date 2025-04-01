package jobs

import (
	"log/slog"
	"time"

	"github.com/senither/dalamud-plugin-listing/state"
)

func StartDeleteExpiredRepositoriesJob(interval time.Duration) {
	tick := time.NewTicker(interval)

	go func() {
		for range tick.C {
			runDelete()
		}
	}()
}

func runDelete() {
	for _, repo := range state.GetRepositories() {
		if repo.RepositoryOrigin.LastUpdatedAt < time.Now().Add(time.Hour*24*3*-1).Unix() {
			var repoUrl string

			if repo.RepoUrl == nil {
				repoUrl = repo.RepositoryOrigin.RepositoryUrl
			} else {
				repoUrl = *repo.RepoUrl
			}

			slog.Info("Deleting expired repository",
				"repository", repo.Name,
				"url", repoUrl,
			)

			state.DeleteRepository(repo)
		}
	}
}
