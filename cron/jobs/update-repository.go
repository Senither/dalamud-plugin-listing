package jobs

import (
	"fmt"
	"time"
)

func UpdateRepository(url string, interval time.Duration, runOnStartup bool) {
	if runOnStartup {
		run(url)
	}

	tick := time.NewTicker(interval)

	go func() {
		for range tick.C {
			run(url)
		}
	}()
}

func run(url string) {
	fmt.Println("Sending request to update repository for: ", url)
}
