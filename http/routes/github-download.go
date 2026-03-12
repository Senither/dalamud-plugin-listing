package routes

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

func DownloadPlugin(c fiber.Ctx) error {
	release := c.Params("*")

	parts := strings.Split(release, "/")
	if len(parts) != 4 {
		return c.Status(fiber.StatusBadRequest).SendString("Bad request, invalid release file format")
	}

	plugin := state.GetInternalPluginByName(parts[0] + "/" + parts[1])
	if plugin == nil || !plugin.Private {
		return c.Status(fiber.StatusNotFound).SendString("Plugin not found")
	}

	releases := state.GetReleaseMetadataByRepositoryName(plugin.Name)
	if releases == nil {
		return c.Status(fiber.StatusNotFound).SendString("No release metadata found for plugin")
	}

	var assetUrl *string = nil
	for _, rel := range releases.Releases {
		if rel.TagName != parts[2] {
			continue
		}

		for _, asset := range rel.Assets {
			if asset.Name == parts[3] {
				assetUrl = &asset.Url
				break
			}
		}
	}

	if assetUrl == nil {
		return c.Status(fiber.StatusNotFound).SendString("Release asset not found")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return c.Status(fiber.StatusInternalServerError).
			SendString("Server misconfigured, missing GITHUB_TOKEN")
	}

	req, err := http.NewRequestWithContext(c, http.MethodGet, *assetUrl, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString("Failed to create upstream request")
	}

	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")
	req.Header.Set("X-Forwarded-For", c.IP())

	slog.Info("Requesting file download for",
		"plugin", plugin.Name,
		"tag", parts[2],
		"asset", parts[3],
		"remote", c.IP(),
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.SendString("Failed to download release asset from GitHub")
		return c.SendStatus(fiber.StatusBadGateway)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))

		return c.Status(fiber.StatusBadGateway).
			SendString("Failed to download release asset from GitHub: " + strings.TrimSpace(string(body)))
	}

	for _, h := range []string{
		"Content-Type",
		"Content-Length",
		"Content-Disposition",
		"Last-Modified",
		"ETag",
		"Cache-Control",
		"Accept-Ranges",
	} {
		if v := resp.Header.Get(h); v != "" {
			c.Set(h, v)
		}
	}

	c.Status(resp.StatusCode)
	_, _ = io.Copy(c.Response().BodyWriter(), resp.Body)

	return nil
}
