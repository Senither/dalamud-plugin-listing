package http

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
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

	app = fiber.New()

	app.Get("/", middleware.RouteSplitter(routes.DashboardHtml, routes.DashboardJson))

	app.Get("/metrics", promhttp.Handler())

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
