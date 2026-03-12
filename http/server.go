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
	app.Get("/download/*", routes.DownloadPlugin)

	app.Get("/plugin/*", middleware.RouteSplitter(routes.PluginHtml, routes.PluginJson))
	app.Get("/", middleware.RouteSplitter(routes.HomepageHtml, routes.HomepageJson))

	app.Use(routes.NotFound)

	app.Listen(addr)
}

func createFiberApp() *fiber.App {
	engine := jet.New("./views", ".jet.html")

	app = fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(favicon.New(favicon.Config{
		File: "./assets/icons/favicon.ico",
		URL:  "/favicon.ico",
	}))

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
