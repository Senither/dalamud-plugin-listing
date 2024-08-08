package main

import (
	"github.com/senither/dalamud-plugin-listing/cron"
	"github.com/senither/dalamud-plugin-listing/http"
)

func main() {
	cron.SetupJobs()
	http.SetupServer()
}
