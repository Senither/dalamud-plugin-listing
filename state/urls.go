package state

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var urls []string

func AddUrl(rawUrl string) {
	if exists(rawUrl) || len(rawUrl) < 4 {
		return
	}

	if !isValidUrl(rawUrl) {
		return
	}

	urls = append(urls, strings.Trim(rawUrl, "\r"))
}

func GetUrls() []string {
	return urls
}

func GetUrlsSize() int {
	return len(GetUrls())
}

func exists(rawUrl string) bool {
	for _, repo := range urls {
		if repo == rawUrl {
			return true
		}
	}

	return false
}

func isValidUrl(rawUrl string) bool {
	url, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		fmt.Println("Error parsing url: ", err)
		return false
	}

	address := net.ParseIP(url.Host)
	if address == nil {
		return strings.Contains(url.Host, ".")
	}

	return true
}
