package clusterdto

import (
	"im.turms/server/internal/domain/common/infra/cluster/discovery"
)

type SettingsDTO struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Settings      map[string]interface{} `json:"settings"`
}

type MemberDTO struct {
	NodeID           string             `json:"nodeId"`
	Zone             string             `json:"zone"`
	Name             string             `json:"name"`
	NodeType         discovery.NodeType `json:"nodeType"`
	IsSeed           bool               `json:"isSeed"`
	IsLeaderEligible bool               `json:"isLeaderEligible"`
	Priority         int                `json:"priority"`
	MemberHost       string             `json:"memberHost"`
	MemberPort       int                `json:"memberPort"`
	AdminAPIAddress  string             `json:"adminApiAddress"`
	WsAddress        string             `json:"wsAddress"`
	TcpAddress       string             `json:"tcpAddress"`
	UdpAddress       string             `json:"udpAddress"`
	IsActive         bool               `json:"isActive"`
	IsHealthy        bool               `json:"isHealthy"`
	IsLeader         bool               `json:"isLeader"`
}

// AddMemberDTO maps to AddMemberDTO in Java.
// Bug fix: Added missing fields (MemberPort, AdminAPIAddress, WsAddress, TcpAddress, UdpAddress, IsActive, IsHealthy, Version)
type AddMemberDTO struct {
	NodeID           string              `json:"nodeId"`
	Zone             string              `json:"zone"`
	Name             string              `json:"name"`
	NodeType         *discovery.NodeType `json:"nodeType"`
	Version          *string             `json:"version"`
	IsSeed           bool                `json:"isSeed"`
	IsLeaderEligible *bool               `json:"isLeaderEligible"`
	Priority         int                 `json:"priority"`
	MemberHost       *string             `json:"memberHost"`
	MemberPort       *int                `json:"memberPort"`
	AdminAPIAddress  *string             `json:"adminApiAddress"`
	WsAddress        *string             `json:"wsAddress"`
	TcpAddress       *string             `json:"tcpAddress"`
	UdpAddress       *string             `json:"udpAddress"`
	IsActive         *bool               `json:"isActive"`
	IsHealthy        *bool               `json:"isHealthy"`
}

type UpdateMemberDTO struct {
	Zone             *string `json:"zone"`
	Name             *string `json:"name"`
	IsSeed           *bool   `json:"isSeed"`
	IsLeaderEligible *bool   `json:"isLeaderEligible"`
	IsActive         *bool   `json:"isActive"`
	Priority         *int    `json:"priority"`
}
