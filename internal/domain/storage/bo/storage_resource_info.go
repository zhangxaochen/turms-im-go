package bo

import "time"

type StorageResourceInfo struct {
	IDNum        *int64
	IDStr        *string
	Name         string
	MediaType    string
	UploaderID   int64
	CreationDate time.Time
}
