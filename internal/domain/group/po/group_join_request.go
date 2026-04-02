package po

import "time"

const CollectionNameGroupJoinRequest = "groupJoinRequest"

// RequestStatus represents the status of a request, mapped to int32 in MongoDB.
// Reusing standard RequestStatus semantics
type RequestStatus int32

const (
	RequestStatusPending RequestStatus = iota
	RequestStatusAccepted
	RequestStatusAcceptedWithoutConfirm
	RequestStatusDeclined
	RequestStatusIgnored
	RequestStatusExpired
	RequestStatusCanceled
)

// GroupJoinRequest represents the group join request entity.
// Sharded by requesterId.
type GroupJoinRequest struct {
	ID           int64         `bson:"_id"`
	Content      string        `bson:"cnt"`
	Status       RequestStatus `bson:"stat"`
	CreationDate time.Time     `bson:"cd"`
	ResponseDate *time.Time    `bson:"rd,omitempty"`
	Reason       *string       `bson:"rsn,omitempty"`
	GroupID      int64         `bson:"gid,omitempty"` // HASH index
	RequesterID  int64         `bson:"rqid"`
	ResponderID  *int64        `bson:"rpid,omitempty"` // HASH index
}
