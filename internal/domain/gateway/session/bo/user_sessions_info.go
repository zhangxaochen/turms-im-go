package bo

import "im.turms/server/pkg/protocol"

// UserSessionInfo represents basic session info returned for admin queries.
type UserSessionInfo struct {
	ID                                   int64               `json:"id"`
	Version                              int32               `json:"version"`
	DeviceType                           protocol.DeviceType `json:"deviceType"`
	DeviceDetails                        map[string]string   `json:"deviceDetails,omitempty"`
	LoginDate                            int64               `json:"loginDate"` // Milliseconds
	LastHeartbeatRequestTimestampMillis int64               `json:"lastHeartbeatRequestTimestampMillis"`
	LastRequestTimestampMillis          int64               `json:"lastRequestTimestampMillis"`
	IsSessionOpen                        bool                `json:"isSessionOpen"`
	Location                             *UserLocation       `json:"location,omitempty"`
	IP                                   []byte              `json:"ip,omitempty"`
}

type UserLocation struct {
	Longitude float32           `json:"longitude"`
	Latitude  float32           `json:"latitude"`
	Timestamp *int64            `json:"timestamp,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

// UserSessionsInfo represents all sessions for a user.
type UserSessionsInfo struct {
	UserID   int64               `json:"userId"`
	Status   protocol.UserStatus `json:"status"`
	Sessions []UserSessionInfo   `json:"sessions"`
}
