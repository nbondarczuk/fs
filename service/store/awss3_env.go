package store

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	. "fs/service/config"	
)

// loadGenericRepositoryEnv Check existence of all AWS S3 specific env variables prefixed
// by a selector whic allows to store many configs like localfile,
// minio, aws s3, etc. and choose one of them by 5prefix of
// the use kind.
// Fields status:
// - host: mandatory, string, address of the bucket server
// - port; optional, string, port to be used if not 80 or 443
// - bucket: mandatory, string, AWS S3 bucket name
// - accessKeyID: mandatory, string, AWS credentials
// - secretAccessKey: mandatory, string, AWS credentials
// - region: mandatory, string, AWS region
// - secure: mandatory, bool, SSL or not
// - presignInt: mandatory, integer, duration of presign url validity in minutes
func LoadGenericRepositoryEnv(prefix string) (err error,
	host, port, bucket, accessKeyID, secretAccessKey, region, secure string,
	presignInt int) {
	var (
		evar   [8]string
		eval   [8]string
	)

	prefix = strings.ToUpper(prefix)
	evar[0] = fmt.Sprintf("%s_HOST", prefix)
	evar[1] = fmt.Sprintf("%s_PORT", prefix)
	evar[2] = fmt.Sprintf("%s_BUCKET_NAME", prefix)
	evar[3] = fmt.Sprintf("%s_ACCESS_KEY_ID", prefix)
	evar[4] = fmt.Sprintf("%s_SECRET_ACCESS_KEY", prefix)
	evar[5] = fmt.Sprintf("%s_REGION", prefix)
	evar[6] = fmt.Sprintf("%s_SECURE", prefix)
	evar[7] = fmt.Sprintf("%s_PRESIGN_DURATION_MIN", prefix)

	// Get the values of variables
	for i := 0; i < 8; i++ {
		eval[i] = os.Getenv(evar[i])
	}

	// Validate the values using their representation
	if eval[0] == "" {
		err = fmt.Errorf("Unknown host in %s: %s", evar[0], eval[0])
	} else if eval[2] == "" {
		err = fmt.Errorf("Unknown bucket name in %s: %s", evar[2], eval[2])
	} else if eval[3] == "" {
		err = fmt.Errorf("Unknown access key id in %s: %s", evar[3], eval[3])
	} else if eval[4] == "" {
		err = fmt.Errorf("Unknown secret access key in %s: %s", evar[4], eval[4])
	} else if eval[5] == "" {
		err = fmt.Errorf("Unknown region in %s: %s", evar[5], eval[5])
	} else if eval[6] == "" {
		err = fmt.Errorf("Unknown secure in %s: %s", evar[6], eval[6])
	} else if eval[7] == "" {
		err = fmt.Errorf("Unknown presign duration in minutes in %s: %s", evar[7], eval[7])
	}
	if err != nil {
		return
	}

	// Convert the values to the final format
	
	host, port, bucket, accessKeyID, secretAccessKey, region, secure =
		eval[0],
		eval[1],
		eval[2],
		eval[3],
		eval[4],
		eval[5],
		eval[6]

	presignInt, err = strconv.Atoi(eval[7])
	if err != nil {
		err = fmt.Errorf("Invalid presign (exp. integer) value in %s: %s - %s",
			evar[7], eval[7], err)
		return
	}

	return
}

// logGenericRepositoryEnv print in the log values of all key env variables
// used to configure the file store repository.
func logGenericRepositoryEnv(host, port, bucket, accessKeyID, secretAccessKey, region, secure string,
	presignInt int) {
	Log.Debug("AWSS3 Repository env:")
	Log.Debug("              Host: " + host)
	Log.Debug("              Port: " + port)
	Log.Debug("            Bucket: " + bucket)
	Log.Debug("     Access Key ID: " + Setup.Anonym(accessKeyID))
	Log.Debug(" Secret Access Key: " + Setup.Anonym(secretAccessKey))
	Log.Debug("            Region: " + region)
	Log.Debug("            Secure: " + secure)	
	Log.Debug("           Presign: " + strconv.FormatInt(int64(presignInt), 10))	
}

// Check existence of all AWS specific env variables
// Fields status:
// - bucket: mandatory, string, AWS S3 bucket name
// - presign: mandatory, integer
func LoadDefaultRepositoryEnv(prefix string) (err error,
	bucket string,
	presignInt int) {
	var (
		evar   [2]string
		eval   [2]string
	)

	prefix = strings.ToUpper(prefix)
	evar[0] = fmt.Sprintf("%s_BUCKET_NAME", prefix)
	evar[1] = fmt.Sprintf("%s_PRESIGN_DURATION_MIN", prefix)

	// Get the values of variables
	for i := 0; i < 2; i++ {
		eval[i] = os.Getenv(evar[i])
	}

	// Validate the values using their representation
	if eval[0] == "" {
		err = fmt.Errorf("Unknown bucket name in %s: %s", evar[0], eval[0])
	} else if eval[1] == "" {
		err = fmt.Errorf("Unknown presign duration in minutes in %s: %s", evar[1], eval[1])
	}
	if err != nil {
		return
	}
	
	// Convert the values to the final format
	
	bucket = eval[0]
	
	presignInt, err = strconv.Atoi(eval[1])
	if err != nil {
		err = fmt.Errorf("Invalid presign (exp. integer) value in %s: %s - %s",
			evar[1], eval[1], err)
	}

	return
}

// logDefaultRepositoryEnv
func logDefaultRepositoryEnv(bucket string,
	presignInt int) {
	Log.Debug("AWSS3 default repository env:")
	Log.Debug(" Bucket: " + bucket)
	Log.Debug("Presign: " + strconv.FormatInt(int64(presignInt), 10))
}
