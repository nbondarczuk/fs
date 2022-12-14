package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	. "fs/service/config"
)

type (
	AppError struct {
		Error      string `json:"error"`
		Message    string `json:"message"`
		HttpStatus int    `json:"status"`
	}

	ErrorResource struct {
		Data AppError `json:"data"`
	}
)

var (
	UrlPathError         = errors.New("URL path decoding error")
	DecoderJsonError     = errors.New("Decoder JSON error")
	EncoderJsonError     = errors.New("Encoder JSON error")
	RepositoryNewError   = errors.New("Repository creation error")
	RepositoryReadError  = errors.New("Repository runtime read error")
	RepositoryWriteError = errors.New("Repository runtime write error")
	RepositoryUseError   = errors.New("Repository runtime use error")	
	PayloadReadError     = errors.New("Payload read error")
	AuthError            = errors.New("Authorisation error")
)

//
// displayAppError showing results in response json stored in header
//
func displayAppError(w http.ResponseWriter, handlerError error, message string, code int) {
	var info string
	if handlerError != nil {
		info = handlerError.Error()
	}

	var ae AppError = AppError{
		Error:      info,
		Message:    message,
		HttpStatus: code,
	}

	Log.Error("Error is: " + fmt.Sprintf("(%d) %s", code, message))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(ErrorResource{Data: ae}); err == nil {
		w.Write(j)
	}
}
