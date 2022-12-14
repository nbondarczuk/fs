package store

import (
	"fmt"

	. "fs/service/config"
)

type EventListener interface {
	ListenForEvents(eventProcessor func(info interface{}) error) error
}

func NewEventListener(kind string) (EventListener, error) {
	Log.Debug("Producing event listener for " + kind)
	
	switch kind[0:5] {
	case "awss3":
		return NewAWSS3EventListener()
	case "minio":
		return NewMinioEventListener()		
	}

	return nil,
		fmt.Errorf("Invalid kind of listener: %s", kind)
}
