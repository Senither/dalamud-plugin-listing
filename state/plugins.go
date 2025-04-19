package state

import (
	"log/slog"
	"strings"
)

var internalPlugins []string

func AddInternalPluginUrl(repoName string) {
	if pluginExists(repoName) || len(repoName) < 4 || !strings.Contains(repoName, "/") {
		return
	}

	internalPlugins = append(internalPlugins, strings.Trim(repoName, "\r"))
}

func GetInternalPlugins() []string {
	return internalPlugins
}

func GetInternalPluginSize() int {
	return len(GetInternalPlugins())
}

func GetSenitherPluginSize() int {
	counter := 0

	for _, repo := range internalPlugins {
		slog.Info("repo", "repo", repo)
		if strings.HasPrefix(strings.ToLower(repo), "senither/") {
			counter++
		}
	}

	return counter
}

func pluginExists(repoName string) bool {
	for _, repo := range internalPlugins {
		if repo == repoName {
			return true
		}
	}

	return false
}
