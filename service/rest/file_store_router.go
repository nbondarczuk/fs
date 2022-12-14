package rest

import (
	"github.com/gorilla/mux"
)

// NewFileStoreRouter creates the router for file service API
func NewFileStoreRouter(r *mux.Router) *mux.Router {
	// Produces presigned URL for POST method. It creates new db record
	// with key info about file and returns POST style URL to create it
	// in the store.
	r.HandleFunc("/api/v1/files",
		CreateFileStore).
		Methods("POST").
		Name("CreateFileStore")

	// Produces presigned URL for GET, PUT, HEAD methods on existing object
	// by numeric primary key. It accesses the db to get the record with key
	// data and then it can access data store bucket where the object is located.
	// Only one object may exist in the db.
	r.HandleFunc("/api/v1/files/{id}/{method:get|put|head}",
		ReadFileStoreAccessById).
		Methods("GET").
		Name("ReadfileStoreAccessById")
	
	// Gets existing db file info by numeric primary key not accessing
	// the data store with object metadata. It must return only one object.
	r.HandleFunc("/api/v1/files/{id:[0-9]+}",
		ReadFileStoreById).
		Methods("GET").
		Name("ReadFileStoreById")

	// Gets existing db file info by fileter on indexed attributes. In theory
	// may return many objects.
	r.HandleFunc("/api/v1/files",
		ReadFileStoreByFilter).
		Methods("GET").
		Name("ReadFileStoreByFilter")
	
	// Updates the status of the db file info by numeric primary key. Only one
	// object can be accessed.
	r.HandleFunc("/api/v1/files/{id:[0-9]+}",
		UpdateFileStoreById).
		Methods("PATCH").
		Name("UpdateFileStoreById")

	return r
}
