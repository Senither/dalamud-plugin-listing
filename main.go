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

	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-shutdownCh
		slog.Info("Shutting down server gracefully")

		slog.Info("Stopping all running jobs...")
		cron.ShutdownJobs()

		runningCh <- struct{}{}

		go func() {
			// This shuts down the application after 3 seconds if it does not exit on its own.
			time.Sleep(3 * time.Second)
			os.Exit(0)
		}()
	}()

	cron.SetupJobs()
	http.SetupServer()

	<-runningCh
}
