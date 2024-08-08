package jobs

import (
	"fmt"
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
			fmt.Println("Deleting expired repository: " + repo.Name + " (" + repo.RepoUrl + ")")

			state.DeleteRepository(repo.RepoUrl)
		}
	}
}
