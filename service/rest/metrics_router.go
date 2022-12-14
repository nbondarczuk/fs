package rest

import (
	"github.com/gorilla/mux"

	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewMetricsRouter
func NewMetricsRouter(r *mux.Router) *mux.Router {
	//r.Handler("GET", "/metrics", promhttp.Handler())
	return r
}

