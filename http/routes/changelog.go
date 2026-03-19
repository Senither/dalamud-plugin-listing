package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

type GitHubReleaseChangelog struct {
	Version   string `json:"version"`
	Changelog string `json:"changelog"`
	CreatedAt string `json:"created_at"`
}

var HasError = fmt.Errorf("Empty error")

func resolveChangelogRequest(c fiber.Ctx) (*state.InternalPlugin, *state.GitHubReleaseContext, string, error) {
	repository, ok := c.Locals("repository").(string)
	if !ok {
		RenderErrorPage(c, fiber.StatusBadRequest, "Bad request", "Bad request, invalid repository name")
		return nil, nil, "", HasError
	}

	parts := strings.Split(repository, "/")
	if len(parts) > 3 || len(parts) < 2 {
		RenderErrorPage(c, fiber.StatusBadRequest, "Bad request", "Bad request, invalid release file format")
		return nil, nil, "", HasError
	}

	plugin := state.GetInternalPluginByName(parts[0] + "/" + parts[1])
	if plugin == nil {
		repositoryPlugin := state.GetRepositoryByAuthorAndInternalName(parts[0], parts[1])
		if repositoryPlugin == nil {
			RenderErrorPage(c, fiber.StatusNotFound, "Plugin not found", "The requested plugin could not be found")
			return nil, nil, "", HasError
		}

		pluginName, _ := strings.CutPrefix(repositoryPlugin.RepositoryOrigin.RepositoryUrl, "https://github.com/")
		plugin = state.GetInternalPluginByName(pluginName)

		if plugin == nil {
			RenderErrorPage(c, fiber.StatusNotFound, "Plugin not found", "The requested plugin could not be found")
			return nil, nil, "", HasError
		}
	}

	releases := state.GetReleaseMetadataByRepositoryName(plugin.Name)
	if releases == nil {
		RenderErrorPage(c, fiber.StatusNotFound, "Release not found", "No release metadata found for plugin")
		return nil, nil, "", HasError
	}

	version := ""
	if len(parts) == 3 {
		version = parts[2]
	}

	return plugin, releases, version, nil
}

func ChangelogHtml(c fiber.Ctx) error {
	plugin, releases, version, err := resolveChangelogRequest(c)
	if err != nil {
		return nil
	}

	if len(version) != 0 {
		c.Redirect().To(fmt.Sprintf("/changelog/%s", plugin.Name))
	}

	var downloadCounter fiber.Map = make(fiber.Map)
	for _, release := range releases.Releases {
		downloadCounter[release.TagName] = 0
		for _, asset := range release.Assets {
			if strings.HasSuffix(asset.Name, ".zip") {
				downloadCounter[release.TagName] = asset.DownloadCount
				break
			}
		}
	}

	return c.Render("changelog", fiber.Map{
		"Plugin":    state.GetRepositoryByGitHubReleaseRepositoryName(plugin.Name),
		"Releases":  releases.Releases,
		"Downloads": downloadCounter,
	}, "layouts/app")
}

func ChangelogJson(c fiber.Ctx) error {
	_, releases, version, err := resolveChangelogRequest(c)
	if err != nil {
		return nil
	}

	if version != "" {
		return renderSingleChangelogEntry(c, releases, version)
	}

	return renderFullChangelogEntries(c, releases)
}

func renderFullChangelogEntries(c fiber.Ctx, releases *state.GitHubReleaseContext) error {
	var changelog []GitHubReleaseChangelog
	for _, release := range releases.Releases {
		changelog = append(changelog, convertGitHubReleaseToChangelogResponse(release))
	}

	return c.JSON(changelog)
}

func renderSingleChangelogEntry(c fiber.Ctx, releases *state.GitHubReleaseContext, version string) error {
	var releaseVersion *state.GitHubPluginRelease
	for _, release := range releases.Releases {
		if release.TagName == version {
			releaseVersion = &release
			break
		}
	}

	if releaseVersion == nil {
		return RenderErrorPage(c, http.StatusNotFound, "Release Not Found", "The requested release version could not be found")
	}

	return c.JSON(convertGitHubReleaseToChangelogResponse(*releaseVersion))
}

func convertGitHubReleaseToChangelogResponse(release state.GitHubPluginRelease) GitHubReleaseChangelog {
	return GitHubReleaseChangelog{
		Version:   release.TagName,
		Changelog: release.Body,
		CreatedAt: release.CreatedAt,
	}
}
