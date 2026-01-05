package http

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/senither/dalamud-plugin-listing/state"
)

func handlePrivatePluginDownload(w http.ResponseWriter, r *http.Request, release string) {
	parts := strings.Split(release, "/")
	if len(parts) != 4 {
		http.Error(w, "Bad request, invalid release file format", http.StatusBadRequest)
		return
	}

	plugin := state.GetInternalPluginByName(parts[0] + "/" + parts[1])
	if plugin == nil || !plugin.Private {
		http.Error(w, "Plugin not found", http.StatusNotFound)
		return
	}

	releases := state.GetReleaseMetadataByRepositoryName(plugin.Name)
	if releases == nil {
		http.Error(w, "No release metadata found for plugin", http.StatusNotFound)
		return
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
		http.Error(w, "Release asset not found", http.StatusNotFound)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		http.Error(w, "Server misconfigured, missing GITHUB_TOKEN", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, *assetUrl, nil)
	if err != nil {
		http.Error(w, "Failed to create upstream request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Dalamud Plugin Listing (https://dalamud-plugins.senither.com/)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to download release asset from GitHub", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		http.Error(w, "Failed to download release asset from GitHub: "+strings.TrimSpace(string(body)), http.StatusBadGateway)
		return
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
			w.Header().Set(h, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
