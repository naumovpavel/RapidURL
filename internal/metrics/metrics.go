package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	TotalRequests    *prometheus.CounterVec
	TotalCacheMisses prometheus.Counter
	ResponseTime     *prometheus.GaugeVec
}

func New(reg prometheus.Registerer) *Metrics {
	metrics := &Metrics{
		TotalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "RapidURL",
				Name:      "total_requests",
				Help:      "Number of total requests.",
			},
			[]string{"url", "method", "status_code"},
		),
		TotalCacheMisses: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "RapidURL",
				Name:      "total_cache_misses",
				Help:      "Number of total cache misses.",
			},
		),
		ResponseTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "RapidURL",
				Name:      "response_time",
				Help:      "Response time for request.",
			},
			[]string{"url", "method", "status_code"},
		),
	}
	reg.MustRegister(metrics.TotalRequests, metrics.ResponseTime, metrics.TotalCacheMisses)
	return metrics
}
