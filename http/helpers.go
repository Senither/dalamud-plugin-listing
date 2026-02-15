package http

import (
	"net"
	"net/http"
	"strings"
)

func GetRequestIP(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); v != "" {
		return v
	}

	if v := strings.TrimSpace(r.Header.Get("True-Client-IP")); v != "" {
		return v
	}

	if v := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); v != "" {
		return strings.TrimSpace(strings.Split(v, ",")[0])
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
