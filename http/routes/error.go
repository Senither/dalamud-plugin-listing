package routes

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func NotFound(c fiber.Ctx) error {
	return RenderErrorPage(c,
		404,
		"Page Not Found",
		"The page or resource you requested does not exist or may have been moved.",
	)
}

func OnlyAcceptsJsonError(c fiber.Ctx) error {
	return RenderErrorPageWithView(c,
		406,
		"Not Acceptable",
		"This endpoint only accepts requests with 'Accept: application/json' header or requests ending with '.json'.",
		"errors/json",
	)
}

func InternalServerError(c fiber.Ctx, err error) error {
	return RenderErrorPage(c, fiber.StatusInternalServerError, "Internal Server Error", err.Error())
}

func RenderErrorPage(c fiber.Ctx, status int, title string, message string) error {
	return RenderErrorPageWithView(c, status, title, message, "errors/error")
}

func RenderErrorPageWithView(c fiber.Ctx, status int, title string, message string, view string) error {
	if isJsonRequest(c) {
		return c.Status(status).JSON(fiber.Map{
			"status": status,
			"reason": message,
			"path":   c.Path(),
		})
	}

	err := c.Status(status).Render(view, fiber.Map{
		"Status":  status,
		"Title":   title,
		"Message": message,
		"Path":    c.Path(),
	}, "layouts/app")

	if err != nil {
		return c.Status(status).SendString(message)
	}

	return nil
}

func isJsonRequest(c fiber.Ctx) bool {
	return c.Get("Accept") == "application/json" ||
		strings.HasPrefix(c.Get("User-Agent"), "Dalamud/") ||
		strings.HasSuffix(c.FullURL(), ".json")
}
