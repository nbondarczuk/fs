package store

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	. "fs/service/config"
)

type AWSS3EventListener struct {
	region           string
	account          string
	notification     string
	objectNamePrefix string
	objectNameSuffix string
	bucketName       string
	session          *session.Session
}

// NewAWSS3eventListener
func NewAWSS3EventListener() (AWSS3EventListener, error) {
	// Load and validate the configuration loaded from config
	err, endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure := loadEventListenerEnv()
	if err != nil {
		return AWSS3EventListener{}, err
	}
	logEventListenerEnv(endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure)

	// Create a session to access the store
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return AWSS3EventListener{}, err
	}
	Log.Debug("Created session: " + fmt.Sprintf("%+v", sess))

	// Store values for later usage in creating objects
	return AWSS3EventListener{
		region:       region,
		account:      account,
		notification: notification,
		bucketName:   bucket,
		session:      sess,
	}, nil
}

// ListenForEvents runs a loop for watching of the object state transitions
func (r AWSS3EventListener) ListenForEvents(eventProcessor func(info interface{}) error) (err error) {
	return
}
