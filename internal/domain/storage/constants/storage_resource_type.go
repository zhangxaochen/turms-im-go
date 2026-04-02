package constants

type StorageResourceType int

const (
	StorageResourceTypeUserProfilePicture StorageResourceType = iota + 1
	StorageResourceTypeGroupProfilePicture
	StorageResourceTypeMessageAttachment
)
