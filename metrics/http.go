package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	routeRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "The total number of HTTP requests made to the server.",
	}, []string{"route"})
)

type RouteMetric string

const (
	HtmlRoute  RouteMetric = "html"
	JsonRoute  RouteMetric = "json"
	ErrorRoute RouteMetric = "error"
)

func IncrementRouteRequestCounter(route RouteMetric) {
	routeRequestCounter.WithLabelValues(string(route)).Inc()
}
