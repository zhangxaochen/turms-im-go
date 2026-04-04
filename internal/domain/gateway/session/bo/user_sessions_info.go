package bo

import "im.turms/server/pkg/protocol"

// UserSessionInfo represents basic session info returned for admin queries.
type UserSessionInfo struct {
	ID         int64               `json:"id"`
	Version    int                 `json:"version"`
	DeviceType protocol.DeviceType `json:"deviceType"`
	LoginDate  int64               `json:"loginDate"` // Milliseconds
	Location   *UserLocation       `json:"location,omitempty"`
}

type UserLocation struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// UserSessionsInfo represents all sessions for a user.
type UserSessionsInfo struct {
	UserID   int64             `json:"userId"`
	Sessions []UserSessionInfo `json:"sessions"`
}
