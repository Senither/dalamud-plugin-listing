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
		"Authors":              state.GetRepositoryAuthors(),
		"Tags":                 state.GetRepositoryTags(),
	}, "layouts/app")
}

func HomepageJson(c fiber.Ctx) error {
	return c.JSON(state.GetRepositories())
}

func RenderPluginListComponent(c fiber.Ctx) error {
	repositories := state.GetRepositories()

	var privatePlugins fiber.Map = make(fiber.Map)
	for _, repo := range repositories {
		privatePlugins[repo.InternalName] = repo.RepositoryOrigin.IsPrivatePlugin != nil &&
			*repo.RepositoryOrigin.IsPrivatePlugin
	}

	sortRepositories(repositories, c.Query("sort"))

	return c.Render("components/plugin-list", fiber.Map{
		"Plugins":         repositories,
		"IsPrivatePlugin": privatePlugins,
	})
}

func sortRepositories(repositories []state.Repository, sortKey string) {
	sort.SliceStable(repositories, func(i, j int) bool {
		left := repositories[i]
		right := repositories[j]

		switch sortKey {
		case "name-desc":
			return left.Name > right.Name
		case "downloads-asc":
			return getValue(left.DownloadCount) > getValue(right.DownloadCount)
		case "downloads-desc":
			return getValue(left.DownloadCount) < getValue(right.DownloadCount)
		case "recently-updated":
			return getLastUpdated(left) > getLastUpdated(right)

		default:
			return left.Name < right.Name
		}
	})
}

func getValue(value any) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	}

	return 0
}

func getLastUpdated(repository state.Repository) int64 {
	var value = int64(0)

	if repository.LastUpdated != nil {
		value = getValue(repository.LastUpdated)
	} else if repository.LastUpdate != nil {
		value = getValue(repository.LastUpdate)
	}

	if value < 1_000_000_000_000 {
		value *= 1000
	}

	return value
}
