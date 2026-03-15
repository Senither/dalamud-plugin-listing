package routes

import (
	"sort"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

func HomepageHtml(c fiber.Ctx) error {
	return c.Render("homepage", fiber.Map{
		"RepositoryCount":      state.GetUrlsSize(),
		"PluginsTotalCount":    state.GetRepositoriesSize(),
		"PluginsInternalCount": state.GetInternalPluginSize(),
		"PluginsSenitherCount": state.GetSenitherPluginSize(),
	}, "layouts/app")
}

func HomepageJson(c fiber.Ctx) error {
	return c.JSON(state.GetRepositories())
}

func RenderPluginListComponent(c fiber.Ctx) error {
	repositories := state.GetRepositories()
	sort.Slice(repositories, func(i, j int) bool {
		return repositories[i].Name < repositories[j].Name
	})

	var privatePlugins fiber.Map = make(fiber.Map)
	for _, repo := range repositories {
		privatePlugins[repo.InternalName] = repo.RepositoryOrigin.IsPrivatePlugin != nil &&
			*repo.RepositoryOrigin.IsPrivatePlugin
	}

	return c.Render("components/plugin-list", fiber.Map{
		"Plugins":         repositories,
		"IsPrivatePlugin": privatePlugins,
	})
}
