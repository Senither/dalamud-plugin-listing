package routes

import (
	"sort"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

type Filter struct {
	Search  string   `query:"search"`
	Tags    []string `query:"tag"`
	Authors []string `query:"author"`
}

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

	var filter Filter
	if err := c.Bind().Query(&filter); err != nil {
		return err
	}

	repositories = filter.bySearch(repositories)
	repositories = filter.byTags(repositories)
	repositories = filter.byAuthors(repositories)
	sortRepositories(repositories, c.Query("sort"))

	return c.Render("components/plugin-list", fiber.Map{
		"Plugins":         repositories,
		"IsPrivatePlugin": privatePlugins,
func (f *Filter) bySearch(repositories []state.Repository) []state.Repository {
	if f.Search == "" {
		return repositories
	}

	query := strings.ToLower(f.Search)

	var filtered []state.Repository
	for _, repo := range repositories {
		if strings.Contains(strings.ToLower(repo.Name), query) ||
			strings.Contains(strings.ToLower(repo.Description), query) ||
			strings.Contains(strings.ToLower(repo.Author), query) {
			filtered = append(filtered, repo)
		}
	}

	return filtered
}

func (f *Filter) byTags(repositories []state.Repository) []state.Repository {
	return filterRepositories(repositories, f.Tags, func(repo state.Repository, value string) bool {
		for _, repoTag := range repo.Tags {
			if strings.ToLower(repoTag) == value {
				return true
			}
		}

		return false
	})
}

func (f *Filter) byAuthors(repositories []state.Repository) []state.Repository {
	return filterRepositories(repositories, f.Authors, func(repo state.Repository, value string) bool {
		for _, repoAuthor := range strings.Split(repo.Author, ",") {
			if strings.TrimSpace(strings.ToLower(repoAuthor)) == value {
				return true
			}
		}

		return false
	})
}

func filterRepositories(
	repositories []state.Repository,
	filters []string,
	match func(repo state.Repository, value string) bool,
) []state.Repository {
	if len(filters) == 0 {
		return repositories
	}

	normalized := make([]string, len(filters))
	for i, v := range filters {
		normalized[i] = strings.ToLower(v)
	}

	var filtered []state.Repository
	for _, repo := range repositories {
		for _, value := range normalized {
			if match(repo, value) {
				filtered = append(filtered, repo)
				break
			}
		}
	}

	return filtered
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
