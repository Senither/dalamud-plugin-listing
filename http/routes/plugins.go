package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

type RepositorySearchCallback func(repo *state.Repository, searchQuery string) bool

func PluginHtml(c fiber.Ctx) error {
	plugin, err := findPluginRepositoryFromContext(c)
	if err != nil {
		return RenderErrorPage(c, 404, "Plugin Not Found", "No plugin was found with the given name.")
	}

	return c.Render("plugin", fiber.Map{
		"Plugin":       plugin,
		"IsInternal":   plugin.RepositoryOrigin.IsInternalPlugin != nil && *plugin.RepositoryOrigin.IsInternalPlugin,
		"IsPrivate":    plugin.RepositoryOrigin.IsPrivatePlugin != nil && *plugin.RepositoryOrigin.IsPrivatePlugin,
		"HasTags":      len(plugin.Tags) > 0,
		"HasChangelog": plugin.Changelog != nil && len(*plugin.Changelog) > 0,
	}, "layouts/app")
}

func PluginJson(c fiber.Ctx) error {
	plugin, err := findPluginRepositoryFromContext(c)
	if err != nil {
		return RenderErrorPage(c, 404, "Plugin Not Found", "No plugin was found with the given name.")
	}

	return c.JSON(plugin)
}

func SearchPluginsByName(c fiber.Ctx) error {
	plugins, err := searchPluginRepositoriesFromContext(c, func(repo *state.Repository, searchQuery string) bool {
		if strings.HasPrefix(strings.ToLower(repo.InternalName), searchQuery) ||
			strings.HasPrefix(strings.ToLower(repo.Name), searchQuery) ||
			strings.Contains(strings.ToLower(repo.Description), searchQuery) {
			return true
		}

		if repo.Punchline != nil && strings.Contains(strings.ToLower(*repo.Punchline), searchQuery) {
			return true
		}

		return false
	})

	if err != nil {
		return RenderErrorPage(c, 404, "No Plugins Found", err.Error())
	}

	return c.JSON(plugins)
}

func SearchPluginsByAuthor(c fiber.Ctx) error {
	plugins, err := searchPluginRepositoriesFromContext(c, func(repo *state.Repository, searchQuery string) bool {
		return strings.Contains(strings.ToLower(repo.Author), searchQuery)
	})

	if err != nil {
		return RenderErrorPage(c, 404, "No Plugins Found", err.Error())
	}

	return c.JSON(plugins)
}

func findPluginRepositoryFromContext(c fiber.Ctx) (*state.Repository, error) {
	name := c.Params("*")
	if before, ok := strings.CutSuffix(name, ".json"); ok {
		name = before
	}

	name = strings.ReplaceAll(name, "%20", " ")

	for _, repo := range state.GetRepositories() {
		if strings.EqualFold(repo.InternalName, name) {
			return &repo, nil
		}
	}

	return nil, fmt.Errorf("No plugin were found with the given name")
}

func searchPluginRepositoriesFromContext(c fiber.Ctx, callback RepositorySearchCallback) ([]state.Repository, error) {
	searchQuery := strings.ToLower(c.Params("*"))
	searchQuery, _ = strings.CutSuffix(searchQuery, ".json")

	var results []state.Repository = make([]state.Repository, 0)

	for _, repo := range state.GetRepositories() {
		if callback(&repo, searchQuery) {
			results = append(results, repo)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("No plugins were found matching the given search query")
	}

	return results, nil
}
