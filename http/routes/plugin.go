package routes

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

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

func findPluginRepositoryFromContext(c fiber.Ctx) (*state.Repository, error) {
	name := c.Params("*")
	slog.Info("Plugin name", "name", name)
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
