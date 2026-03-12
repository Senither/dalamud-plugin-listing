package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/metrics"
)

func RouteSplitter(htmlRoute fiber.Handler, jsonRoute fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		if isJsonRequest(c) {
			metrics.IncrementRouteRequestCounter(metrics.JsonRoute)
			return jsonRoute(c)
		}

		metrics.IncrementRouteRequestCounter(metrics.HtmlRoute)
		return htmlRoute(c)
	}
}

func isJsonRequest(c fiber.Ctx) bool {
	return c.Get("Accept") == "application/json" ||
		strings.HasPrefix(c.Get("User-Agent"), "Dalamud/") ||
		strings.HasSuffix(c.FullURL(), ".json")
}
