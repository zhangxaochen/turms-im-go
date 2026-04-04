package po

import "time"

// Meeting represents a meeting entity in MongoDB (meeting collection).
// @MappedFrom Meeting
type Meeting struct {
	ID           int64      `bson:"_id"`
	CreatorID    int64      `bson:"cid"`
	UserID       *int64     `bson:"uid,omitempty"`
	GroupID      *int64     `bson:"gid,omitempty"`
	CreationDate time.Time  `bson:"cd"`
	StartDate    time.Time  `bson:"sd"`
	CancelDate   *time.Time `bson:"cad,omitempty"`
	EndDate      *time.Time `bson:"ed,omitempty"`
	Name         *string    `bson:"n,omitempty"`
	Intro        *string    `bson:"intro,omitempty"`
	Password     *string    `bson:"pw,omitempty"`
}

const MeetingCollectionName = "meeting"
