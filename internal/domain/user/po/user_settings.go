package po

import "time"

const CollectionNameUserSettings = "userSettings"

// UserSettings represents application-level settings related to a user.
// MongoDB Collection: userSettings
type UserSettings struct {
	UserID          int64                  `bson:"_id"` // userId
	Settings        map[string]interface{} `bson:"s"`   // settings
	LastUpdatedDate time.Time              `bson:"lud"` // lastUpdatedDate
}
