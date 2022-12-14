package rest

import (
	"fs/service/db"
)

type (
	// input: for Create entity of S3 bucket file allocation object
	FileStoreRequestResource struct {
		CheckSum string `json:"check_sum,omitempty"`
		Name     string `json:"name,omitempty"`
		Size     int64  `json:"size,omitempty"`
		Status   string `json:"status,omitempty"`
	}

	// output: S3 bucket file allocation object id with an access URL address from POST
	FileStoreReplyResource struct {
		Status bool           `json:"status"`
		Count  int64          `json:"count"`
		Data   []db.FileStore `json:"data,omitempt"`
	}
)
