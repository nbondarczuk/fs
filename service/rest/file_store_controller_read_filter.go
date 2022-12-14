package rest

import (
	"encoding/json"
	"net/http"

	. "fs/service/config"
	"fs/service/db"
)

// Gets existing file allocation object by filter on tid, did
func ReadFileStoreByFilter(w http.ResponseWriter, r *http.Request) {
	tenant := r.FormValue("tenant")
	if tenant == "" {
		displayAppError(w, UrlPathError,
			"Missing mandatory parameter tenant",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got parameter tenant: " + tenant)

	device := r.FormValue("device") 
	if device == "" {
		displayAppError(w, UrlPathError,
			"Missing mandatory paramter device",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got parameter device: " + device)

	dbrep, err := db.NewRepository()
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer dbrep.Close()

	fss, count, err := dbrep.ReadByFilter(tenant, device)
	if err != nil {
		displayAppError(w, RepositoryReadError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	var reply = FileStoreReplyResource{
		Status: true,
		Count:  count,
		Data:   fss,
	}

	jstr, err := json.Marshal(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while marshalling data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	Log.Debug("Reply: " + string(jstr))
	
	writeResponseWithJson(w, http.StatusOK, jstr)
}
