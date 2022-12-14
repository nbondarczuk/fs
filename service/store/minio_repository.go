package store

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	. "fs/service/config"
)

// Stores necessary data to perform object operations
type MinioRepository struct {
	bucket   string
	object   string
	client   *minio.Client
	presign  int
	region   string
	filename string
}

// FullObjectName produces canonical name of the object to be used in cross-identification
// saving theoriginal data source kind full name, not only prefix of it. It can identify
// the location of the data bucket, ie. AWS, local instance, etc.
func (r MinioRepository) FullObjectName() string {
	return Setup.UseFileStore + ":" + r.bucket + ":" + r.object
}

// NewMinioRepository creates a object with client connection to the datastore.
// It uses the config for AWS key env variables to make the connection.
func NewMinioRepository(tid, did, name string) (MinioRepository, error) {
	// Load, validate, convert the configuration loaded from config
	err, endpoint, port, bucket, accessKeyID, secretAccessKey, region, presign, secure := loadRepositoryEnv()
	if err != nil {
		return MinioRepository{},
			fmt.Errorf("Error collecting env values: %s", err)
	}
	object := objectKeyName(tid, did, name)
	logRepositoryEnv(endpoint, port, bucket, object, accessKeyID, secretAccessKey, region, presign, secure)
	// Do the conversion in 2nd step as some of the value may be malformed
	presignVal, err := strconv.Atoi(presign)
	if err != nil {
		return MinioRepository{},
			fmt.Errorf("Invalid presign (exp. integer) value: %s - %s", secure, err)
	}
	secureVal, err := strconv.ParseBool(secure)
	if err != nil {
		return MinioRepository{},
			fmt.Errorf("Invalid secure (exp. bool) value: %s - %s", secure, err)
	}

	// Create a client to access the store via host and specific port if specified
	// Local MINIO is running on 9000 and not the default 80 or 443 ports so to get connected
	// the specific port has to be added to the endpoint address.
	if port != "" {
		endpoint = endpoint + ":" + port
	}
	Log.Debug("Connecting endpoint: " + endpoint)
	
	client, err := minio.New(endpoint,
		&minio.Options{
			//client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: secureVal,
			Region: region,
		})
	if err != nil {
		return MinioRepository{}, err
	}
	if Setup.LogLevel == "debug" {
		client.TraceOn(os.Stdout)
	}
	Log.Debug("Created client: " + fmt.Sprintf("%+v", client))

	// Store values for later usage in creating store bucket objects
	return MinioRepository{
		bucket:   bucket,
		object:   object,
		client:   client,
		presign:  presignVal,
		region:   region,
		filename: name,
	}, nil
}

// FileObjectPresignedPostURL produces presigned URL to be accessed by the clients
// with a seq of name/content form data to be used in curl (for example) call
// with -F or --form option for each of the mappings. The object is not created
// as it will be done with resigned POST operation.
func (r MinioRepository) FileObjectPresignedURL(cksum string, size int64) (string, map[string]string, error) {
	// The bucket must exist, if not, then create it
	err := r.AssureBucketExist()
	if err != nil {
		return "", nil, fmt.Errorf("Error checking bucket; %s", err)
	}

	// Policy with upload restrictions
	policy := minio.NewPostPolicy()
	policy.SetBucket(r.bucket)
	policy.SetKey(r.object)
	policy.SetExpires(time.Now().UTC().AddDate(0, 0, r.presign)) // expires in N day(s)
	policy.SetContentLengthRange(0, size)

	// Get the POST ready URL with form key/value object
	url, formData, err := r.client.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		return "", nil, fmt.Errorf("Error getting presigned url: %s", err)
	}
	Log.Debug("Got presigned URL: " + url.String() + " - " + fmt.Sprintf("%+v", formData))

	return url.String(), formData, nil
}

// FileObjectPresignedGetURL provides url for GET mode acceess to the existing object.
// It mey be used to implement Read operation on the object. The precondition is
// that the object exists.
func (r MinioRepository) FileObjectPresignedGetURL(cksum string, size int64) (string, error) {
	// The bucket and object must exist, if not, then error as POST had to be done
	err := r.AssureBucketObjectExist()
	if err != nil {
		return "", fmt.Errorf("Error checking bucket object: %s", err)
	}

	// Additional response header overrides supports:
	// response-expires, response-content-type, response-cache-control, response-content-disposition
	params := make(url.Values)
	disposition := fmt.Sprintf("attachment; filename=\"%s\"", r.filename)
	params.Set("response-content-disposition", disposition)

	// Produce presigned URL
	url, err := r.client.PresignedGetObject(context.Background(),
		r.bucket,
		r.object,
		time.Second * 24 * 60 * 60 * time.Duration(r.presign),
		params)

	return url.String(), nil
}

// FileObjectPresignedHeadURL provides url for HEAD mode acceess to the existing object.
// It mey be used to implement Read operation on the object. The precondition is
// that the object exists.
func (r MinioRepository) FileObjectPresignedHeadURL(cksum string, size int64) (string, error) {
	// The bucket and object must exist, if not, then error as POST had to be done
	err := r.AssureBucketObjectExist()
	if err != nil {
		return "", fmt.Errorf("Error checking bucket object: %s", err)
	}

	// Additional response header overrides supports:
	// response-expires, response-content-type, response-cache-control, response-content-disposition
	params := make(url.Values)
	disposition := fmt.Sprintf("attachment; filename=\"%s\"", r.filename)
	params.Set("response-content-disposition", disposition)

	// Produce presigned URL
	url, err := r.client.PresignedHeadObject(context.Background(),
		r.bucket,
		r.object,
		time.Second * 24 * 60 * 60 * time.Duration(r.presign),
		params)

	return url.String(), nil
}

// FileObjectPresignedPutURL provides url for GET mode acceess to the existing object.
// It may be used to implement Update operation on the object. The precondition is
// that the object exists.
func (r MinioRepository) FileObjectPresignedPutURL(cksum string, size int64) (string,  error) {
	// The bucket object must exist, if not, then error as POST had to be done
	err := r.AssureBucketObjectExist()
	if err != nil {
		return "", fmt.Errorf("Error checking bucket object: %s", err)
	}

	// Produce presigned URL
	url, err := r.client.PresignedPutObject(context.Background(),
		r.bucket,
		r.object,
		time.Second * 24 * 60 * 60 * time.Duration(r.presign))

	return url.String(), nil
}

// AssureBucketExist checks bucket existence and creates bucket if it does not exist
func (r MinioRepository) AssureBucketExist() error {
	Log.Debug("Checking bucket: " + r.bucket)

	found, err := r.client.BucketExists(context.Background(), r.bucket)
	if err != nil {
		return err
	}
	Log.Debug("Bucket status checked: " + r.bucket)
	if !found {
		Log.Debug("Bucket does not exist: " + r.bucket)
		err = r.client.MakeBucket(context.Background(),
			r.bucket,
			minio.MakeBucketOptions{
				Region:        r.region,
				ObjectLocking: false,
			})
		if err != nil {
			return err
		}
		Log.Debug("Created bucket: " + r.bucket)
	}

	return nil
}

// AssureBucketObjectExist checks bucket and object existence
func (r MinioRepository) AssureBucketObjectExist() error {
	Log.Debug("Checking bucket object: " + r.object)

	found, err := r.client.BucketExists(context.Background(), r.bucket)
	if err != nil {
		return err
	}
	Log.Debug("Bucket status checked: " + r.bucket)
	if !found {
		return fmt.Errorf("Bucket does not exist: " + r.bucket)
	}

	// Returns error if object does not exist
	info, err := r.client.StatObject(context.Background(),
		r.bucket,
		r.object,
		minio.StatObjectOptions{})
	if err != nil {
		return err
	}
	Log.Debug("Info bucket object: " + fmt.Sprintf("%+v", info))
	
	return nil
}
