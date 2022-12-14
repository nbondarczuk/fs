package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Status struct {
	RequestsNo    int
	CacheMissesNo int
	TenentFileOps map[string]int
	DeviceFeleOps map[string]int
}

var (
	RequestsCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "requests_counter",
		Help:        "Number of requests received",
		ConstLabels: prometheus.Labels{"version": "1"},
	})
)

func init() {
	prometheus.MustRegister(RequestsCounter)
}
