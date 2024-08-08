package state

import "strings"

var urls []string

func AddUrl(url string) {
	if exists(url) || len(url) < 4 {
		return
	}

	urls = append(urls, strings.Trim(url, "\r"))
}

func GetUrls() []string {
	return urls
}

func exists(url string) bool {
	for _, repo := range urls {
		if repo == url {
			return true
		}
	}

	return false
}
