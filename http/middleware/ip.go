package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RequestIP(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	if v := strings.TrimSpace(c.Get("CF-Connecting-IP")); v != "" {
		return output.WriteString(v)
	}

	if v := strings.TrimSpace(c.Get("True-Client-IP")); v != "" {
		return output.WriteString(v)
	}

	if v := strings.TrimSpace(c.Get("X-Forwarded-For")); v != "" {
		return output.WriteString(strings.TrimSpace(strings.Split(v, ",")[0]))
	}

	return output.WriteString(c.IP())
}
