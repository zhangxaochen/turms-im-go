package po

import "time"

const ConversationSettingsCollectionName = "conversationSettings"

// ConversationSettings represents the settings for a conversation.
type ConversationSettings struct {
	// ID is the composite key of ownerId and targetId.
	ID ConversationSettingsKey `bson:"_id"`

	// Settings is the map of settings.
	// BSON key "s" stands for "settings".
	Settings map[string]any `bson:"s,omitempty"`

	// LastUpdatedDate is the date when the settings were last updated.
	// BSON key "lud" stands for "lastUpdatedDate".
	LastUpdatedDate time.Time `bson:"lud"`
}

type ConversationSettingsKey struct {
	OwnerId  int64 `bson:"oid"`
	TargetId int64 `bson:"tid"`
}

const (
	ConversationSettingsFieldIdOwnerId       = "_id.oid"
	ConversationSettingsFieldSettings         = "s"
	ConversationSettingsFieldLastUpdatedDate = "lud"
)
