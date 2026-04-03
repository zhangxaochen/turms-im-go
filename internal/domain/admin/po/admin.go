package po

import "time"

const CollectionNameAdmin = "admin"

type Admin struct {
	ID               int64     `bson:"_id"`
	LoginName        string    `bson:"ln"`
	Password         []byte    `bson:"pw"`
	DisplayName      string    `bson:"n"`
	RoleIDs          []int64   `bson:"rid,omitempty"`
	RegistrationDate time.Time `bson:"rd,omitempty"`
}

// Fields mapping for Admin
const (
	AdminFieldID               = "_id"
	AdminFieldLoginName        = "ln"
	AdminFieldPassword         = "pw"
	AdminFieldDisplayName      = "n"
	AdminFieldRoleIDs          = "rid"
	AdminFieldRegistrationDate = "rd"
)
