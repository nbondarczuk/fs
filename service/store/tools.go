package store

import (
	"fmt"
	"os"
	"strings"
	
	. "fs/service/config"	
)

// objectKey provides uniform way of naming s3 bucket objects
func objectKeyName(tid, did, name string) string {
	return fmt.Sprintf("%s-%s-%s", tid, did, name)
}

// Check existence of all AWS S3 specific env variables prefixed
// by a selector whic allows to store many configs like localfile,
// minio, aws s3, etc. and choose one of them by 5prefix of
// the use kind.
// Fields status:
// - endpoint: mandatory, string, address of the bucket server
// - port; optional, string, port to be used if not 80 or 443
// - bucket: mandatory, string, AWS S3 bucket name
// - accessKeyID: mandatory, string, AWS credentials
// - secretAccessKey: mandatory, string, AWS credentials
// - region: mandatory, string, AWS region
// - presign: mandatory, integer, duration of presign url validity in days
// - secure: mandatory, bool, SSL or not
func loadRepositoryEnv() (err error,
	endpoint, port, bucket, accessKeyID, secretAccessKey, region, presign, secure string) {
	var (
		prefix string = strings.ToUpper(Setup.UseFileStore)
		arg1   string = fmt.Sprintf("%s_ENDPOINT", prefix)
		arg2   string = fmt.Sprintf("%s_PORT", prefix)
		arg3   string = fmt.Sprintf("%s_BUCKET_NAME", prefix)
		arg4   string = fmt.Sprintf("%s_ACCESS_KEY_ID", prefix)
		arg5   string = fmt.Sprintf("%s_SECRET_ACCESS_KEY", prefix)
		arg6   string = fmt.Sprintf("%s_REGION", prefix)
		arg7   string = fmt.Sprintf("%s_PRESIGN_DURATION_DAYS", prefix)
		arg8   string = fmt.Sprintf("%s_SECURE", prefix)
	)
	endpoint, port, bucket, accessKeyID, secretAccessKey, region, presign, secure = os.Getenv(arg1),
		os.Getenv(arg2),
		os.Getenv(arg3),
		os.Getenv(arg4),
		os.Getenv(arg5),
		os.Getenv(arg6),
		os.Getenv(arg7),
		os.Getenv(arg8)
	if endpoint == "" {
		err = fmt.Errorf("Unknown endpoint from %s", arg1)
	} else if bucket == "" {
		err = fmt.Errorf("Unknown bucket name from %", arg3)
	} else if accessKeyID == "" {
		err = fmt.Errorf("Unknown access key id from %", arg4)
	} else if secretAccessKey == "" {
		err = fmt.Errorf("Unknown secret access key from %", arg5)
	} else if region == "" {
		err = fmt.Errorf("Unknown region from %", arg6)
	} else if region == "" {
		err = fmt.Errorf("Unknown presign duration in days from %", arg7)
	} else if secure == "" {
		err = fmt.Errorf("Unknown secure value %", arg7)
	}

	return
}

// logRepositoryEnv print in the log values of all key env variables
// used to configure the file store repository.
func logRepositoryEnv(endpoint, port, bucket, object, accessKeyID, secretAccessKey, region, presign, secure string) {
	Log.Debug("AWSS3 Repository env:")
	Log.Debug("            Endpoint: " + endpoint)
	Log.Debug("                Port: " + port)
	Log.Debug("              Bucket: " + bucket)
	Log.Debug("              Object: " + object)
	Log.Debug("       Access Key ID: " + accessKeyID)
	Log.Debug("   Secret Access Key: " + secretAccessKey)
	Log.Debug("              Region: " + region)
	Log.Debug("             Presign: " + presign)
	Log.Debug("              Secure: " + secure)
}

// Check existence of all AWS specific env variables
// Fields status:
// - endpoint: mandatory, string, address of the bucket server
// - port; optional, string, port to be used if not 80 or 443
// - bucket: mandatory, string, AWS S3 bucket name
// - accessKeyID: mandatory, string, AWS credentials
// - secretAccessKey: mandatory, string, AWS credentials
// - account: mandatory, string, AWS credentials
// - region: mandatory, string, AWS region
// - notification: mandatory, string, name of notification topic
// - secure: mandatory, bool, SSL or not 
func loadEventListenerEnv() (err error, endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure string) {
	var (
		prefix string = strings.ToUpper(Setup.UseFileStore)
		arg1   string = fmt.Sprintf("%s_ENDPOINT", prefix)
		arg2   string = fmt.Sprintf("%s_PORT", prefix)
		arg3   string = fmt.Sprintf("%s_BUCKET_NAME", prefix)
		arg4   string = fmt.Sprintf("%s_ACCESS_KEY_ID", prefix)
		arg5   string = fmt.Sprintf("%s_SECRET_ACCESS_KEY", prefix)
		arg6   string = fmt.Sprintf("%s_ACCOUNT_ID", prefix)		
		arg7   string = fmt.Sprintf("%s_REGION", prefix)
		arg8   string = fmt.Sprintf("%s_NOTIFICATION_NAME", prefix)
		arg9   string = fmt.Sprintf("%s_SECURE", prefix)		
	)
	
	endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure = os.Getenv(arg1),
		os.Getenv(arg2),
		os.Getenv(arg3),
		os.Getenv(arg4),
		os.Getenv(arg5),
		os.Getenv(arg6),
		os.Getenv(arg7),
		os.Getenv(arg8),
		os.Getenv(arg9)		
	if endpoint == "" {
		err = fmt.Errorf("Unknown endpoint")
	} else if bucket == "" {
		err = fmt.Errorf("Unknown bucket name")
	} else if accessKeyID == "" {
		err = fmt.Errorf("Unknown access key id")
	} else if secretAccessKey == "" {
		err = fmt.Errorf("Unknown secret access key")
	} else if account == "" {
		err = fmt.Errorf("Unknown account")
	} else if region == "" {
		err = fmt.Errorf("Unknown region")
	} else if notification == "" {
		err = fmt.Errorf("Unknown notification name")
	} else if secure == "" {
		err = fmt.Errorf("Unknown secure value")
	}

	return
}

// logEventListenerEnv
func logEventListenerEnv(endpoint, port, bucket, accessKeyID, secretAccessKey, account, region, notification, secure string) {
	Log.Debug("AWSS3 Event Listener env:")
	Log.Debug("         Endpoint: " + endpoint)
	Log.Debug("             Port: " + port)	
	Log.Debug("           Bucket: " + bucket)
	Log.Debug("    Access Key ID: " + accessKeyID)
	Log.Debug("Secret Access Key: " + secretAccessKey)
	Log.Debug("          Account: " + account)
	Log.Debug("           Region: " + region)
	Log.Debug("     Notification: " + notification)
	Log.Debug("           Secure: " + secure)	
}
