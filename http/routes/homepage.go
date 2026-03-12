package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

func HomepageHtml(c fiber.Ctx) error {
	return c.Render("homepage", fiber.Map{
		"RepositoryCount":      state.GetUrlsSize(),
		"PluginsTotalCount":    state.GetRepositoriesSize(),
		"PluginsInternalCount": state.GetInternalPluginSize(),
		"PluginsSenitherCount": state.GetSenitherPluginSize(),
		"Plugins":              state.GetRepositories(),
	})
}

func HomepageJson(c fiber.Ctx) error {
	return c.JSON(state.GetRepositories())
}
