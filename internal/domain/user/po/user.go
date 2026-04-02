package po

import "time"

// User represents a user entity in MongoDB (user collection).
type User struct {
	ID                int64      `bson:"_id"`
	Password          string     `bson:"pw,omitempty"`
	Name              string     `bson:"n,omitempty"`
	Intro             string     `bson:"intro,omitempty"`
	ProfilePicture    string     `bson:"pp,omitempty"`
	ProfileAccess     int32      `bson:"pas"` // enum for strategy
	PermissionGroupID int64      `bson:"pgid"`
	RegistrationDate  time.Time  `bson:"rd"`
	DeletionDate      *time.Time `bson:"dd,omitempty"`
	LastUpdatedDate   *time.Time `bson:"lud,omitempty"`
	IsActive          bool       `bson:"act,omitempty"`

	UserDefinedAttributes map[string]any `bson:"user_defined_attributes,omitempty"`
}
