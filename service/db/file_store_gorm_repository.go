package db

import (
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fs/service/config"
	. "fs/service/config"
)

// GORM file allocation handle
type FileStoreRepositoryGORM struct {
	be     BackendPostgres
	gormdb *gorm.DB
}

// NewFileStoreRepository creates a handle for domain operations
// on a bucket allocations using gorm
func NewRepository() (FileStoreRepositoryGORM, error) {
	// Backend connection
	credentials, err := NewBackendCredentialsPostgres()
	if err != nil {
		return FileStoreRepositoryGORM{}, err
	}
	Log.Debug("Validated Postgres credentials")
	backend, err := NewBackendPostgres(credentials)
	if err != nil {
		return FileStoreRepositoryGORM{}, err
	}
	Log.Debug("Created backend Postgres")
	err = backend.Ping()
	if err != nil {
		return FileStoreRepositoryGORM{}, err
	}
	Log.Debug("Pinged backend Postgres DB")

	// GORM creation
	var gormConfig gorm.Config
	if config.Setup.LogLevel == "debug" {
		gormConfig = gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		gormConfig = gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: backend.Sqldb,
	}), &gormConfig)
	if err != nil {
		return FileStoreRepositoryGORM{}, err
	}
	Log.Debug("Opened GORM on backend DB")
	gormdb.AutoMigrate(&FileStore{})
	Log.Debug("Migrated object FileStore")

	// Store results for use
	return FileStoreRepositoryGORM{
		backend,
		gormdb,
	}, nil
}

// Create new record allocating new ID
func (r FileStoreRepositoryGORM) Create(tid, did, name, cksum string, size int64) ([]FileStore, int64, error) {
	fa := FileStore{
		TenantID:     tid,
		DeviceID:     did,
		CheckSumType: Setup.CheckSumType,
		CheckSum:     cksum,
		Name:         name,
		Size:         size,
		Status:       "N",
	}
	result := r.gormdb.Create(&fa)
	fas := make([]FileStore, 1)
	fas[0] = fa

	return fas, result.RowsAffected, result.Error
}

// ReadById - primary key read by ID of the record
func (r FileStoreRepositoryGORM) ReadById(id string) ([]FileStore, int64, error) {
	// Check the parameters
	idpkval, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	// Read the entity by primary key access
	var fa FileStore
	result := r.gormdb.First(&fa, idpkval)
	fas := make([]FileStore, 1)
	fas[0] = fa

	return fas, result.RowsAffected, result.Error
}

// ReadByFilter - fiter by indexed (hopefully) attributes read
func (r FileStoreRepositoryGORM) ReadByFilter(tid, did string) ([]FileStore, int64, error) {
	// Read the entity by fuilter on non-key attributes
	var fas []FileStore
	result := r.gormdb.Where("tenant_id = ? AND device_id = ?", tid, did).Find(&fas)

	return fas, result.RowsAffected, result.Error
}

// UpdateById - primary key update by ID of the record
func (r FileStoreRepositoryGORM) UpdateById(id, status, cksum string, size int64) ([]FileStore, int64, error) {
	// Check the parameters
	idpkval, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, 0, err
	}
	if status == "" {
		return nil, 0, fmt.Errorf("Invalid status, empty string")
	}

	// Do the GORM style update of selected parameters plus mandatory status
	var fa FileStore
	fa.ID = uint(idpkval)
	result := r.gormdb.Model(&fa).Updates(FileStore{
		CheckSumType: Setup.CheckSumType,
		CheckSum:     cksum,
		Size:         size,
		Status:       status,
	})
	fas := make([]FileStore, 1)
	fas[0] = fa

	return fas, result.RowsAffected, result.Error
}

// SetBucketLocation
func (r FileStoreRepositoryGORM) SetBucketLocation(id uint, bucket, location string) error {
	var fa FileStore
	fa.ID = id
	result := r.gormdb.Model(&fa).Updates(FileStore{
		Location: location,
		Bucket:   bucket,
	})

	return result.Error
}

// Close
func (r FileStoreRepositoryGORM) Close() {
	r.be.Close()
}
