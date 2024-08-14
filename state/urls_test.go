package state

import "testing"

func init() {
	urls = []string{}
}

func teardown() {
	urls = []string{}
}

func TestAddValidUrl(t *testing.T) {
	defer teardown()

	url := "https://example.com"

	AddUrl(url)

	if len(urls) != 1 {
		t.Errorf("Expected 1 url, got %d", len(urls))
	}

	if urls[0] != url {
		t.Errorf("Expected url to be %s, got %s", url, urls[0])
	}
}

func TestAddInvalidUrl(t *testing.T) {
	defer teardown()

	url := "some-invalid-url"

	AddUrl(url)

	if len(urls) != 0 {
		t.Errorf("Expected 0 url, got %d", len(urls))
	}
}
