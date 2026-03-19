package http

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/responsetime"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/jet/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/senither/dalamud-plugin-listing/http/middleware"
	"github.com/senither/dalamud-plugin-listing/http/routes"
)

var app *fiber.App

func SetupServer() {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	app = createFiberApp()

	app.Get("/assets/*", static.New("./assets"))
	app.Get("/metrics", promhttp.Handler())

	app.Post("/webhook/github-release", routes.GitHubReleaseWebhook)

	app.Get("/download/*", middleware.ParseRepositoryParam, routes.DownloadPlugin)
	app.Get("/plugin/*", middleware.ParseRepositoryParam, middleware.RouteSplitter(routes.PluginHtml, routes.PluginJson))
	app.Get("/plugins/*", middleware.ParseRepositoryParam, middleware.RouteSplitter(routes.OnlyAcceptsJsonError, routes.SearchPluginsByName))
	app.Get("/authors/*", middleware.ParseRepositoryParam, middleware.RouteSplitter(routes.OnlyAcceptsJsonError, routes.SearchPluginsByAuthor))
	app.Get("/changelog/*", middleware.ParseRepositoryParam, middleware.RouteSplitter(routes.ChangelogHtml, routes.ChangelogJson))

	app.Get("/", middleware.RouteSplitter(routes.HomepageHtml, routes.HomepageJson))

	hx := app.Group("/hx")

	hx.Get("/plugins", routes.RenderPluginListComponent)

	app.Use(routes.NotFound)

	app.Listen(addr)
}

func createFiberApp() *fiber.App {
	engine := jet.New("./views", ".jet.html")

	app = fiber.New(fiber.Config{
		AppName:      "Dalamud Plugin List",
		Views:        engine,
		ErrorHandler: routes.InternalServerError,
	})

	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET")
		c.Set("Access-Control-Allow-Headers", "accept,content-type")
		c.Set("Access-Control-Max-Age", "300")

		c.Set("Server", "Dalamud Plugin List")
		c.Set("X-Powered-By", "Nerd Rage and Caffeine")
		c.Set("X-Accepts", "text/html,application/json")

		return c.Next()
	})

	app.Use(favicon.New(favicon.Config{
		File: "./assets/icons/favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Use(logger.New(logger.Config{
		CustomTags: map[string]logger.LogFunc{"realIP": middleware.RequestIP},
		Format:     "${time} | ${status} | ${latency} | ${realIP} | ${method} | ${path} | ${error}\n",
	}))

	app.Use(responsetime.New())

	engine.Templates.AddGlobal("StyleHash", generateStyleHashId())

	return app
}

func generateStyleHashId() string {
	file, err := os.ReadFile("./assets/styles.css")
	if err != nil {
		slog.Error("Failed to read styles.css", "err", err)
		return ""
	}

	return fmt.Sprintf("%x", sha256.Sum256(file))
}

func ShutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Graceful shutdown error: %v", "err", err)
	}

	slog.Info("Gracefully shutdown the server")
}
