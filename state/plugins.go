package state

import (
	"strings"
)

type InternalPlugin struct {
	Name    string
	Private bool
}

var internalPlugins []InternalPlugin

func AddInternalPluginUrl(repoName string) {
	var private = false

	if strings.HasPrefix(repoName, "P:") {
		private = true
		repoName = strings.TrimPrefix(repoName, "P:")
	}

	if pluginExists(repoName) || len(repoName) < 4 || !strings.Contains(repoName, "/") {
		return
	}

	internalPlugins = append(internalPlugins, InternalPlugin{
		Name:    repoName,
		Private: private,
	})
}

func GetInternalPluginByName(repoName string) *InternalPlugin {
	for _, repo := range internalPlugins {
		if strings.EqualFold(repo.Name, repoName) {
			return &repo
		}
	}

	return nil
}

func GetInternalPlugins() []InternalPlugin {
	return internalPlugins
}

func GetInternalPluginSize() int {
	return len(GetInternalPlugins())
}

func GetSenitherPluginSize() int {
	counter := 0

	for _, repo := range internalPlugins {
		if strings.HasPrefix(strings.ToLower(repo.Name), "senither/") {
			counter++
		}
	}

	return counter
}

func pluginExists(repoName string) bool {
	for _, repo := range internalPlugins {
		if strings.EqualFold(repo.Name, repoName) {
			return true
		}
	}

	return false
}
