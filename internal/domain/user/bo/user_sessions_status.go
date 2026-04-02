package bo

import (
	"im.turms/server/pkg/protocol"
)

type UserSessionsStatus struct {
	UserID                        int64
	UserStatus                    protocol.UserStatus
	OnlineDeviceTypeToSessionInfo map[protocol.DeviceType]UserDeviceSessionInfo
}

func (s UserSessionsStatus) GetNodeIDIfActive(deviceType protocol.DeviceType) string {
	if info, ok := s.OnlineDeviceTypeToSessionInfo[deviceType]; ok && info.IsActive {
		return info.NodeID
	}
	return ""
}

type UserDeviceSessionInfo struct {
	NodeID                    string
	HeartbeatTimestampSeconds int64
	IsActive                  bool
}

type UserStatusFieldType int

const (
	UserStatusFieldTypeUserStatus UserStatusFieldType = iota
	UserStatusFieldTypeDeviceTypeToNodeID
	UserStatusFieldTypeNodeIDToHeartbeatTimestamp
)

type UserStatusField struct {
	Type  UserStatusFieldType
	Value interface{}
}
