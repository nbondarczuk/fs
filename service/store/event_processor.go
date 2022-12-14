package store

import (
	"fmt"

	"github.com/minio/minio-go/v7/pkg/notification"

	. "fs/service/config"
	"fs/service/db"
)

// EventProcessor mapps the json from the event queue into DB object reference
// and it changes the status of if according to the event info. The most
// important status change is that onject was uploaded or downloaded.
// The mapping goes like this:
//   s3:ObjectCreated:*  -> (C)reated
// doing update on the initial record created in state = (N)ew
// The access to the object may cause change to cksum or size
// so the values must be put in sync.
func EventProcessor(info notification.Info) error {
	Log.Info("Processing events: " + fmt.Sprintf("%d", len(info.Records)))

	// New handle for DB access with object record refer to object by index
	r, err := db.NewRepository()
	if err != nil {
		return fmt.Errorf("Error creating db repository: %s", err)
	}
	defer r.Close()
	
	// For each event received within this info
	for _, event := range info.Records {
		Log.Info("Processing event " + fmt.Sprintf("%+v", info))

		// Extract key data from the event
		fullObjectName := "aws:s3:" + event.S3.Bucket.Name + ":" + event.S3.Object.Key
		status := mapEventNameToStatus(event.EventName)

		// Get the db record in sync with file store
		count, err := r.UpdateByFullObjectName(fullObjectName, status)
		if err != nil {
			return fmt.Errorf("Error updating db repository: %s", err)
		} else if count != 1 {
			Log.Error("Invalid number of records updated")
		}
	}
	
	return nil
}

// mapEventNameToStatus converts event code to new status codes according
// to the table:
//   s3:ObjectCreated:*  -> (C)reated
func mapEventNameToStatus(name string) string {
	return "C" // tbd
}
