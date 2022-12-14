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

// Creates new record for an object in file store and saves its reference
// in the db
func CreateFileStore(w http.ResponseWriter, r *http.Request) {
	tenant := r.FormValue("tenant")
	if tenant == "" {
		displayAppError(w, UrlPathError,
			"Missing mandatory parameter tenant",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got paramter tenant: " + tenant)

	device := r.FormValue("device")
	if device == "" {
		displayAppError(w, UrlPathError,
			"Missing mandatory parameter device",
			http.StatusBadRequest)
		return
	}
	Log.Debug("Got paramter device: " + device)

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
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer dbrep.Close()

	// The ID is allocated by a unique seq to be used as part of the name of the object
	fss, count, err := dbrep.Create(tenant, device, request.Name, request.CheckSum, request.Size)
	if err != nil {
		displayAppError(w, RepositoryWriteError,
			"Error while creating entity - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	id := fss[0].ID
	Log.Debug("Created db record with ID: " + fmt.Sprintf("%d", id))

	// As the bucket name may be allocated on the fly in case of more then 1 buckets
	fsrep, err := store.NewRepository(Setup.UseFileStore, tenant, device, request.Name, id)
	if err != nil {
		displayAppError(w, RepositoryNewError,
			"Error while creating repository - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	location, bucket := fsrep.BucketLocation()

	// The initially created db entry must be updated and the bucket registered
	err = dbrep.SetBucketLocation(id, bucket, location)
	if err != nil {
		displayAppError(w, RepositoryWriteError,
			"Error while linking object bucket to location of object - "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Now is the time to get a presigned URL pointing to an object to be created
	// by PUT operation by the client
	url, err := fsrep.FileObjectPresignedPutURL(request.CheckSum, request.Size)
	if err != nil {
		displayAppError(w, RepositoryUseError,
			"Error while allocating file object url - "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	Log.Debug("Created client PUT presigned url: " + url)

	// The presigned url and seq of name/content form data is not stored in the DB
	// but it must be returned to the client.
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
