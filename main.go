package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/senither/dalamud-plugin-listing/cron"
	"github.com/senither/dalamud-plugin-listing/http"
)

func main() {
	runningCh := make(chan struct{}, 1)
	shutdownCh := make(chan os.Signal, 1)

	if os.Getenv("APP_URL") == "" {
		slog.Warn("APP_URL environment variable is not set, defaulting to http://localhost:8080/")
		os.Setenv("APP_URL", "http://localhost:8080/")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		slog.Warn("GITHUB_TOKEN environment variable is not set, some features may not work properly")
	}

	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-shutdownCh
		slog.Info("Shutting down server gracefully")

		slog.Info("Stopping all running jobs...")
		cron.ShutdownJobs()

		slog.Info("Shutting down the HTTP server...")
		http.ShutdownServer()

		runningCh <- struct{}{}

		go func() {
			// This shuts down the application after 10 seconds if it does not exit on its own.
			time.Sleep(10 * time.Second)
			slog.Info("Forcing the application to exit after 10 seconds")
			os.Exit(0)
		}()
	}()

	cron.SetupJobs()
	http.SetupServer()

	<-runningCh
}
