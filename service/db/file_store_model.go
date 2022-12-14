package db

import (
	"gorm.io/gorm"
)

// S3 bucket file object handle with an access URL address.
// The status state transition is: (W)aiting -> (U)ploaded, (D)ownloaded, (E)xpired
// It is done with PATCH so file size or cksum may be changed as well in transit.
type FileStore struct {
	gorm.Model

	Location     string `json:"location,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	TenantID     string `gorm:"index:idx_tenant_device" json:"tenent_id,omitempty"`
	DeviceID     string `gorm:"index:idx_tenant_device" json:"device_id,omitempty"`
	Name         string `json:"name,omitempty"`
	CheckSumType string `json:"check_sum_type,omitempty"`
	CheckSum     string `json:"check_sum,omitempty"`
	Size         int64  `json:"size,omitempty"`
	Status       string `json:"status,omitempty"`
	URL          string `gorm:"-" json:"url,omitempty"`
}
