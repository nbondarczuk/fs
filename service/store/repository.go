package store

import (
	"fmt"

	. "fs/service/config"
)

// Repository is an interface method to be used by the API to access the data store
type Repository interface {
	BucketLocation() (string, string)	
	FileObjectPresignedGetURL(cksum string, size int64) (string, error)
	FileObjectPresignedPutURL(cksum string, size int64) (string, error)
	FileObjectPresignedHeadURL(cksum string, size int64) (string, error)		
}

// NewRepository is a dispatch for specific repositories implemented with MINIO
// or some other technology like AWS S3 API, Google or Azure, supported
// by MINIO or not but having (possibly) different structure of access crdentials
// and/or config options.
func NewRepository(kind, tid, did, name string, id uint) (Repository, error) {
	Log.Debug("Producing repository of kind: " + kind +
		" for object: " + objectKeyName(id, tid, did,  name))

	if kind == "awss3" { // AWSS3 uses default AWS config
		return NewAWSS3RepositoryDefault(tid, did, name, id)
	} else if kind[0:5] == "awss3" { // Other generic full setup
		return NewAWSS3RepositoryGeneric(tid, did, name, id)
	}

	return nil,
		fmt.Errorf("Invalid kind of repository: %s", kind)
}
