package store

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/notification"

	. "fs/service/config"
)

type MinioEventListener struct {
	region           string
	account          string
	notification     string
	objectNamePrefix string
	objectNameSuffix string
	bucketName       string
	client           *minio.Client
}

// NewMinioeventListener
func NewMinioEventListener() (MinioEventListener, error) {
	// Load and validate the configuration loaded from config
	err, endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure := loadEventListenerEnv()
	if err != nil {
		return MinioEventListener{}, err
	}
	logEventListenerEnv(endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure)
	secureVal, err := strconv.ParseBool(secure)
	if err != nil {
		return MinioEventListener{},
			fmt.Errorf("Invalid secure (exp. bool) value: %s - %s", secure, err)
	}

	// Create a client to access the store
	if port != "" {
		endpoint = endpoint + ":" + port
	}
	Log.Debug("Connecting endpoint: " + endpoint)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secureVal,
		Region: region,
	})
	if err != nil {
		return MinioEventListener{}, err
	}
	if Setup.LogLevel == "debug" {
		client.TraceOn(os.Stderr)
	}
	Log.Debug("Created client: " + fmt.Sprintf("%+v", client))

	// Store values for later usage in creating objects
	return MinioEventListener{
		region:       region,
		account:      account,
		notification: notification,
		bucketName:   bucket,
		client:       client,
	}, nil
}

// AWS SQS configuration created in the AWS and no changes to the receiver
func (r MinioEventListener) configEventListener() error {
	// For AWS SQS notification method is configured
	queueArn := notification.NewArn("aws", "sqs", r.region, r.account, r.notification)

	// Add list of events of some interest and filter for the name
	queueConfig := notification.NewConfig(queueArn)
	queueConfig.AddFilterPrefix(r.objectNamePrefix)
	queueConfig.AddFilterSuffix(r.objectNameSuffix)
	queueConfig.AddEvents(notification.ObjectCreatedAll)
	config := notification.Configuration{}
	config.AddQueue(queueConfig)

	// Confure the the bucket to notify about the specific events
	err := r.client.SetBucketNotification(context.Background(),
		r.bucketName,
		config)
	if err != nil {
		return fmt.Errorf("Unable to set the bucket notification: ", err)
	}

	return nil
}

// ListenForEvents runs a loop for watching of the object state transitions
func (r MinioEventListener) ListenForEvents(eventProcessor func(info interface{}) error) error {
	return nil
}
