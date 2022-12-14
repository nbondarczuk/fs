package rest

import (
	"encoding/json"
	"net/http"
	
	. "fs/service/config"
	"fs/service/db"
)

// Updates the metricsus of the existing file allocation object by primary key
func UpdateFileStoreById(w http.ResponseWriter, r *http.Request) {
	id, err := pathVariableStr(r, "id", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable id",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got path variable id: " + id)

	payload, err := readPayload(r)
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to read payload of the request",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got text payload: " + string(payload))

	var request FileStoreRequestResource
	err = json.Unmarshal(payload, &request)
	if err != nil {
		displayAppError(w, PayloadReadError,
			"Unable to decode json payload of the request",
			http.StatusBadRequest)
		return
	}
	
	dbrep, err := db.NewRepository()
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating db repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer dbrep.Close()

	fss, count, err := dbrep.UpdateById(id, request.Status, request.CheckSum, request.Size)
	if err != nil {
		displayAppError(w, RepositoryWriteError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	if count == 0 {
		displayAppError(w, RepositoryReadError,
			"Error no item updated by id - "+err.Error(),
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
