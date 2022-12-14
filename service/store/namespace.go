package store

import (
	"fmt"

	. "fs/service/config"
)

const (
	objectKeyNameSeparator string = "-"
)

// objectKeyName provides uniform way of naming s3 bucket objects.
// The prefix of the object name is a module of its internal id thus
// distributing the objects over the bucket namespace evenly as
// the module value is a function of its unique id.
func objectKeyName(id uint, tid, did, name string) string {
	return fmt.Sprintf("%d%s%s%s%s%s%s",
		id % Setup.ObjectNamespaceModule, objectKeyNameSeparator,
		tid, objectKeyNameSeparator,
		did, objectKeyNameSeparator,
		name)
}
