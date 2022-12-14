package rest

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router for API collecting all sub-routers.
// They handle route groups like client admin or system routes.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r = NewSystemRouter(r)
	r = NewMetricsRouter(r)	
	r = NewFileStoreRouter(r)

	return r
}
