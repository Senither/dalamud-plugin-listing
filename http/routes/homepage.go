package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/senither/dalamud-plugin-listing/state"
)

func DashboardHtml(c fiber.Ctx) error {
	return c.SendFile("./views/index.html")
}

func DashboardJson(c fiber.Ctx) error {
	return c.JSON(state.GetRepositories())
}
