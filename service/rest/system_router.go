package rest

import (
	"github.com/gorilla/mux"
)

// NewSystemRouter
func NewSystemRouter(r *mux.Router) *mux.Router {
	r.HandleFunc("/system/health",
		ReadSystemHealth).
		Methods("GET").
		Name("read-system-health")

	r.HandleFunc("/system/alive",
		ReadSystemAlive).
		Methods("GET").
		Name("read-system-alive")

	r.HandleFunc("/system/stat",
		ReadSystemStat).
		Methods("GET").
		Name("read-system-stat")

	r.HandleFunc("/system/version",
		ReadSystemVersion).
		Methods("GET").
		Name("read-system-version")

	return r
}
