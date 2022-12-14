package rest

import (
	"encoding/json"
	"net/http"

	. "fs/service/config"
	"fs/service/db"
)

// Gets existing file allocation object by primary key
func ReadFileStoreById(w http.ResponseWriter, r *http.Request) {
	id, err := pathVariableStr(r, "id", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable id",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got path variable id: " + id)

	dbrep, err := db.NewRepository()
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating db repository - "+err.Error(),
			http.StatusBadRequest)
		return
	}
	defer dbrep.Close()

	fss, count, err := dbrep.ReadById(id)
	if err != nil {
		displayAppError(w, RepositoryReadError,
			"Error while reading from db repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	if count == 0 {
		displayAppError(w, RepositoryReadError,
			"Error no item found by id - "+err.Error(),
			http.StatusNotFound)
		return
	}

	if count > 1 {
		displayAppError(w, RepositoryReadError,
			"Error no item found by id - "+err.Error(),
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
