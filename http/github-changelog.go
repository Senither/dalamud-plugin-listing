package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/senither/dalamud-plugin-listing/state"
)

type GitHubReleaseChangelog struct {
	Version   string `json:"version"`
	Changelog string `json:"changelog"`
	CreatedAt string `json:"created_at"`
}

func handleInternalPluginChangelog(w http.ResponseWriter, r *http.Request, release string) {
	parts := strings.Split(release, "/")
	if len(parts) > 3 || len(parts) < 2 {
		http.Error(w, "Bad request, invalid release file format", http.StatusBadRequest)
		return
	}

	plugin := state.GetInternalPluginByName(parts[0] + "/" + parts[1])
	if plugin == nil {
		http.Error(w, "Plugin not found", http.StatusNotFound)
		return
	}

	releases := state.GetReleaseMetadataByRepositoryName(plugin.Name)
	if releases == nil {
		http.Error(w, "No release metadata found for plugin", http.StatusNotFound)
		return
	}

	if len(parts) == 3 {
		renderInternalPluginSingleReleaseChangelog(w, releases, parts[2])
		return
	}

	var changelog []GitHubReleaseChangelog
	for _, release := range releases.Releases {
		changelog = append(changelog, convertGitHubReleaseToChangelogResponse(release))
	}

	content, err := json.Marshal(changelog)
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(content))
}

func renderInternalPluginSingleReleaseChangelog(w http.ResponseWriter, releases *state.GitHubReleaseContext, version string) {
	var releaseVersion *state.GitHubPluginRelease
	for _, release := range releases.Releases {
		if release.TagName == version {
			releaseVersion = &release
			break
		}
	}

	if releaseVersion == nil {
		http.Error(w, "Release version not found", http.StatusNotFound)
		return
	}

	content, err := json.Marshal(convertGitHubReleaseToChangelogResponse(*releaseVersion))
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(content))
}

func convertGitHubReleaseToChangelogResponse(release state.GitHubPluginRelease) GitHubReleaseChangelog {
	return GitHubReleaseChangelog{
		Version:   release.TagName,
		Changelog: release.Body,
		CreatedAt: release.CreatedAt,
	}
}
