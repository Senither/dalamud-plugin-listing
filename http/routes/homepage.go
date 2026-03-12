package routes

import (
	"sort"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

func HomepageHtml(c fiber.Ctx) error {
	repositories := state.GetRepositories()
	sort.Slice(repositories, func(i, j int) bool {
		return repositories[i].Name < repositories[j].Name
	})

	return c.Render("homepage", fiber.Map{
		"RepositoryCount":      state.GetUrlsSize(),
		"PluginsTotalCount":    state.GetRepositoriesSize(),
		"PluginsInternalCount": state.GetInternalPluginSize(),
		"PluginsSenitherCount": state.GetSenitherPluginSize(),
		"Plugins":              repositories,
	}, "layouts/app")
}

func HomepageJson(c fiber.Ctx) error {
	return c.JSON(state.GetRepositories())
}
