package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// Initial count.
	currentCount                  = 0
	currentCountGetAll            = 0
	currentCountGetById           = 0
	currentCountGetByIdAndVersion = 0
	currentCountLabels            = 0
	currentCountCreate            = 0
	currentCountExtend            = 0
	currentCountDelete            = 0

	// The Prometheus metric that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total",
			Help: "Total number of http hits.",
		},
	)

	httpHitsGetAll = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_get_all",
			Help: "Total number of http hits.",
		},
	)

	httpHitsGetById = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_get_by_id",
			Help: "Total number of http hits.",
		},
	)

	httpHitsGetByIdVersion = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_get_by_id_version",
			Help: "Total number of http hits.",
		},
	)

	httpHitsLabels = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_by_labels",
			Help: "Total number of http hits.",
		},
	)

	httpHitsCreate = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_create",
			Help: "Total number of http hits.",
		},
	)

	httpHitsExtend = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_extend",
			Help: "Total number of http hits.",
		},
	)

	httpHitsDelete = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total_delete",
			Help: "Total number of http hits.",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		httpHits,
		httpHitsGetAll,
		httpHitsGetById,
		httpHitsGetByIdVersion,
		httpHitsLabels,
		httpHitsCreate,
		httpHitsExtend,
		httpHitsDelete,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		f(w, r) // original function call
	}
}

func countGetAll(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsGetAll.Inc()
		f(w, r) // original function call
	}
}

func countGetById(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsGetById.Inc()
		f(w, r) // original function call
	}
}

func countGetByIdVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsGetByIdVersion.Inc()
		f(w, r) // original function call
	}
}

func countGetByLabels(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsLabels.Inc()
		f(w, r) // original function call
	}
}

func countGetByCreate(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsCreate.Inc()
		f(w, r) // original function call
	}
}

func countGetByExtend(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsExtend.Inc()
		f(w, r) // original function call
	}
}

func countGetByDelete(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHitsDelete.Inc()
		f(w, r) // original function call
	}
}
