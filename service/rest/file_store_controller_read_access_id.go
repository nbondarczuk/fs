package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	. "fs/service/config"
	"fs/service/db"
	"fs/service/store"
)

// ReadFileStoreAccessById gets a presigned GET/PUT/HEAD URL to perform an operation
// on the existing file store resource assuming that the resource does exist,
// ie. was previously created by POST using presigned URL.
func ReadFileStoreAccessById(w http.ResponseWriter, r *http.Request) {
	id, err := pathVariableStr(r, "id", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable id",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got path variable id: " + id)

	method, err := pathVariableStr(r, "method", true)
	if err != nil {
		displayAppError(w, UrlPathError,
			"Missing mandatory url path variable method",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got path variable method: " + method)

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

	fsrep, err := store.NewRepository(Setup.UseFileStore, fss[0].TenantID, fss[0].DeviceID, fss[0].Name, fss[0].ID)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// TBD - What about a map of func?
	var url string
	switch method {
	case "head":
		url, err = fsrep.FileObjectPresignedHeadURL(fss[0].CheckSum, fss[0].Size)
	case "get":
		url, err = fsrep.FileObjectPresignedGetURL(fss[0].CheckSum, fss[0].Size)
	case "put":
		url, err = fsrep.FileObjectPresignedPutURL(fss[0].CheckSum, fss[0].Size)
	default:
		err = fmt.Errorf("Invalid access method: %s, expecting: head, get, put", method)
	}
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while producing file object presigned url - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// The presigned url of the existing object both in the db as in the file store
	fss[0].URL = url

	var reply = FileStoreReplyResource{
		Status: true,
		Count:  count,
		Data:   fss,
	}

	// Eascape chars shall not be replaces by unicodes as the standard MArshall does
	var writer bytes.Buffer
	enc := json.NewEncoder(&writer)
	enc.SetEscapeHTML(false)
	err = enc.Encode(&reply)
	if err != nil {
		displayAppError(w, EncoderJsonError,
			"An error while encoding response data - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	jstr := writer.Bytes()
	Log.Debug("Reply: " + string(jstr))

	writeResponseWithJson(w, http.StatusOK, jstr)
}
