package store_test

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	. "fs/service/config"
	"fs/service/store"
)

// Mock setup
func setTestConfig(t *testing.T) {
	s, err := NewSetupValueSet([]byte(""))
	if err != nil {
		t.Fatalf("Error loading valid setup: %s", err)
	}
	if s == nil {
		t.Fatalf("Empty Setup")
	}

	s.UseFileStore = "testing"
	
	Setup = s
}

func TestDefaultRepositoryEnv(t *testing.T) {
	setTestConfig(t)

	t.Run("detect empty bucket name", func(t *testing.T) {
		os.Setenv("TESTING_BUCKET_NAME", "")
		os.Setenv("TESTING_PRESIGN_DURATION_MIN", "345")
		
		err, _, _ := store.LoadDefaultRepositoryEnv(Setup.UseFileStore)
		if err == nil {
			t.Fatalf("Error loading default repository env")
		}
	})

	t.Run("detect invalid integer format", func(t *testing.T) {
		os.Setenv("TESTING_BUCKET_NAME", "XXX")
		os.Setenv("TESTING_PRESIGN_DURATION_MIN", "xxx")
		
		err, _, _ := store.LoadDefaultRepositoryEnv(Setup.UseFileStore)
		if err == nil {
			t.Fatalf("Error loading default repository env")
		}
	})
	
	t.Run("gets all env vars needed", func(t *testing.T) {
		os.Setenv("TESTING_BUCKET_NAME", "123")
		os.Setenv("TESTING_PRESIGN_DURATION_MIN", "345")
		
		err, bucket, presign := store.LoadDefaultRepositoryEnv(Setup.UseFileStore)
		if err != nil {
			t.Fatalf("Error loading default repository env: %s", err.Error())
		}

		assert.Equal(t, bucket, "123")
		assert.Equal(t, presign, 345)
	})
}

func TestGenericRepositoryEnv(t *testing.T) {
	setTestConfig(t)

	t.Run("detect empty bucket name", func(t *testing.T) {
		os.Setenv("TESTING_HOST", "testhost")
		os.Setenv("TESTING_PORT", "testport")
		os.Setenv("TESTING_BUCKET_NAME", "")
		os.Setenv("TESTING_ACCESS_KEY_ID", "testaccesskeyid")
		os.Setenv("TESTING_SECRET_ACCESS_KEY", "testsecretaccesskey")
		os.Setenv("TESTING_REGION", "testregion")
		os.Setenv("TESTING_SECURE", "testsecure")
		os.Setenv("TESTING_PRESIGN_DURATION_MIN", "345")
		
		err, _, _,_, _, _, _, _, _ := store.LoadGenericRepositoryEnv(Setup.UseFileStore)
		if err == nil {
			t.Fatalf("Error loading default repository env")
		}
	})

	t.Run("detect invalid integer format", func(t *testing.T) {
		os.Setenv("TESTING_HOST", "testhost")
		os.Setenv("TESTING_PORT", "testport")
		os.Setenv("TESTING_BUCKET_NAME", "testbucket")
		os.Setenv("TESTING_ACCESS_KEY_ID", "testaccesskeyid")
		os.Setenv("TESTING_SECRET_ACCESS_KEY", "testsecretaccesskey")
		os.Setenv("TESTING_REGION", "testregion")
		os.Setenv("TESTING_SECURE", "testsecure")
		os.Setenv("TESTING_PRESIGN_DURATION_MIN", "xxx")
		
		err, _, _,_, _, _, _, _, _  := store.LoadGenericRepositoryEnv(Setup.UseFileStore)
		if err == nil {
			t.Fatalf("Error loading default repository env")
		}
	})
	
	t.Run("gets all env vars needed", func(t *testing.T) {
		os.Setenv("TESTING_HOST", "testhost")
		os.Setenv("TESTING_PORT", "testport")
		os.Setenv("TESTING_BUCKET_NAME", "testbucket")
		os.Setenv("TESTING_ACCESS_KEY_ID", "testaccesskeyid")
		os.Setenv("TESTING_SECRET_ACCESS_KEY", "testsecretaccesskey")
		os.Setenv("TESTING_REGION", "testregion")
		os.Setenv("TESTING_SECURE", "testsecure")
		os.Setenv("TESTING_PRESIGN_DURATION_MIN", "123")
		
		err, host, port, bucket, accessKeyID, secretAccessKey, region, secure, presign := store.LoadGenericRepositoryEnv(Setup.UseFileStore)
		if err != nil {
			t.Fatalf("Error loading default repository env: %s", err.Error())
		}

		assert.Equal(t, host, "testhost")
		assert.Equal(t, port, "testport")
		assert.Equal(t, bucket, "testbucket")
		assert.Equal(t, accessKeyID, "testaccesskeyid")
		assert.Equal(t, secretAccessKey, "testsecretaccesskey")
		assert.Equal(t, region, "testregion")
		assert.Equal(t, secure, "testsecure")
		assert.Equal(t, presign, 123)
	})
}
