package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	. "fs/service/config"
)

// Stores necessary data to perform object operations
type AWSS3Repository struct {
	object   string
	bucket   string
	presign  time.Duration
	filename string
	session  *session.Session
	service  *s3.S3
}

// LocationOfBucket
func (r AWSS3Repository) BucketLocation() (string, string) {
	return Setup.UseFileStore, r.bucket
}

// NewAWSS3RepositoryDefault creates a object with client connection to the datastore.
// It uses default config AWS key env variables to make the connection.
func NewAWSS3RepositoryDefault(tid, did, name string, id uint) (AWSS3Repository, error) {
	object := objectKeyName(id, tid, did, name)

	// Load, validate, convert the config
	err, bucket, presign := LoadDefaultRepositoryEnv(Setup.UseFileStore)
	if err != nil {
		return AWSS3Repository{},
			fmt.Errorf("Error collecting env values: %s", err)
	}
	logDefaultRepositoryEnv(bucket, presign)

	// Determine the log level of AWS connection
	var logLevel *aws.LogLevelType
	if Setup.LogLevel == "debug" {
		logLevel = aws.LogLevel(aws.LogDebugWithHTTPBody)
	}

	// Create an AWS session to access the store
	sess, err := session.NewSession(&aws.Config{
		LogLevel: logLevel,
	})
	if err != nil {
		return AWSS3Repository{}, err
	}
	Log.Debug("Created session: " + fmt.Sprintf("%+v", *sess))

	// Create S3 service client
	svc := s3.New(sess, aws.NewConfig().WithLogLevel(*logLevel))

	// Store values for later usage in creating store bucket objects
	return AWSS3Repository{
		object:   object,
		bucket:   bucket,
		presign:  time.Duration(presign) * time.Minute,
		filename: name,
		session:  sess,
		service:  svc,
	}, nil
}

// NewAWSS3RepositoryGeneric creates a object with client connection to the datastore.
// It uses the config for AWS key env variables to make the connection.
// It is used for connection to the local testing store butt it is compatible
// with the AWS S3 connection.
func NewAWSS3RepositoryGeneric(tid, did, name string, id uint) (AWSS3Repository, error) {
	object := objectKeyName(id, tid, did, name)

	// Load, validate, convert the config
	err, host, port, bucket, accessKeyID, secretAccessKey, region, secure, presign := LoadGenericRepositoryEnv(Setup.UseFileStore)
	if err != nil {
		return AWSS3Repository{},
			fmt.Errorf("Error collecting env values: %s", err)
	}
	logGenericRepositoryEnv(host, port, bucket, accessKeyID, secretAccessKey, region, secure, presign)

	// All details from config file so that we can use minio for local testing
	// but AWS may be accessed this way as well

	creds := credentials.NewStaticCredentials(accessKeyID,
		secretAccessKey,
		"")

	if port == "" {
		port = "80" // Well known default port if not defined in the config
	}

	var protocol string
	if "true" == strings.ToLower(secure) {
		protocol = "https"
	} else {
		protocol = "http"
	}

	endpoint := fmt.Sprintf("%s://%s:%s", protocol, host, port)

	// Determine the log level of AWS connection
	var logLevel *aws.LogLevelType
	if Setup.LogLevel == "debug" {
		logLevel = aws.LogLevel(aws.LogDebugWithHTTPBody)
	}

	// Create an AWS session to access the store
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
		LogLevel:         logLevel,
	})
	if err != nil {
		return AWSS3Repository{}, err
	}
	Log.Debug("Created session: " + fmt.Sprintf("%+v", sess))

	// Create S3 service client
	svc := s3.New(sess, aws.NewConfig().WithLogLevel(*logLevel))

	// Store values for later usage in creating store bucket objects
	return AWSS3Repository{
		object:   object,
		bucket:   bucket,
		presign:  time.Duration(presign) * time.Minute,
		filename: name,
		session:  sess,
		service:  svc,
	}, nil
}

// FileObjectPresignedGetURL provides url for GET mode acceess to the existing object.
// It mey be used to implement Read operation on the object. The precondition is
// that the object exists.
func (r AWSS3Repository) FileObjectPresignedGetURL(cksum string, size int64) (string, error) {
	// GET method to be used
	req, _ := r.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket:         aws.String(r.bucket),
		Key:            aws.String(r.object),
	})

	// Get presinged URL of create empty object
	url, err := req.Presign(r.presign)
	if err != nil {
		return "", fmt.Errorf("Failed to sign request: %s", err)
	}
	Log.Debug("Created AWS GET presigned request url: " + url)

	return url, nil
}

// FileObjectPresignedHeadURL provides url for HEAD mode acceess to the existing object.
// It mey be used to implement Read operation on the object. The precondition is
// that the object exists.
func (r AWSS3Repository) FileObjectPresignedHeadURL(cksum string, size int64) (string, error) {
	// HEAD method to be used
	req, _ := r.service.HeadObjectRequest(&s3.HeadObjectInput{
		Bucket:         aws.String(r.bucket),
		Key:            aws.String(r.object),
	})

	// Get presinged URL of create empty object
	url, err := req.Presign(r.presign)
	if err != nil {
		return "", fmt.Errorf("Failed to sign request: %s", err)
	}
	Log.Debug("Created AWS HEAD presigned request url: " + url)

	return url, nil
}

// FileObjectPresignedPutURL provides url for GET mode acceess to the existing object.
// It may be used to implement Update operation on the object. The precondition is
// that the object exists.
func (r AWSS3Repository) FileObjectPresignedPutURL(cksum string, size int64) (string, error) {
	// PUT method to be used
	req, _ := r.service.PutObjectRequest(&s3.PutObjectInput{
		Bucket:         aws.String(r.bucket),
		Key:            aws.String(r.object),
		ChecksumSHA256: aws.String(cksum),
		ContentLength:  aws.Int64(size),
	})

	// Get presinged URL of create empty object
	url, err := req.Presign(r.presign)
	if err != nil {
		return "", fmt.Errorf("Failed to sign request: %s", err)
	}
	Log.Debug("Created AWS PUT presigned request url: " + url)

	return url, nil
}
