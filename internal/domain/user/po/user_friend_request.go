package po

import "time"

const CollectionNameUserFriendRequest = "userFriendRequest"

// RequestStatus represents the status of a friend request.
type RequestStatus int32

const (
	RequestStatusPending RequestStatus = iota
	RequestStatusAccepted
	RequestStatusDeclined
	RequestStatusIgnored
	RequestStatusExpired
	RequestStatusCanceled
)

// ResponseAction represents the response action of a friend request.
type ResponseAction int32

const (
	ResponseActionAccept ResponseAction = iota
	ResponseActionDecline
	ResponseActionIgnore
)

// UserFriendRequest represents the user friend request entity.
type UserFriendRequest struct {
	ID           int64         `bson:"_id"`
	Content      string        `bson:"c"`
	Status       RequestStatus `bson:"s"`
	Reason       *string       `bson:"r,omitempty"`
	CreationDate time.Time     `bson:"cd"`
	ResponseDate *time.Time    `bson:"rd,omitempty"`
	RequesterID  int64         `bson:"rqid"`
	RecipientID  int64         `bson:"rcid"`
}
