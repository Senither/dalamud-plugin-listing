package http

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/jet/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/senither/dalamud-plugin-listing/http/middleware"
	"github.com/senither/dalamud-plugin-listing/http/routes"
)

var app *fiber.App

func SetupServer(views embed.FS) {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	app = fiber.New(fiber.Config{
		Views: jet.New("./views", ".jet"),
	})

	app.Get("/assets/*", static.New("./assets"))
	app.Get("/metrics", promhttp.Handler())

	app.Get("/", middleware.RouteSplitter(routes.HomepageHtml, routes.HomepageJson))

	app.Listen(addr)
}

func ShutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Graceful shutdown error: %v", "err", err)
	}

	slog.Info("Gracefully shutdown the server")
}
