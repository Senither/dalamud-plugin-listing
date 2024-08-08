package jobs

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
