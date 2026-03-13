package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func ParseRepositoryParam(c fiber.Ctx) error {
	name := c.Params("*")
	if after, ok := strings.CutSuffix(name, ".json"); ok {
		name = after
	}

	c.Locals("repository", strings.ReplaceAll(name, "%20", " "))

	return c.Next()
}
