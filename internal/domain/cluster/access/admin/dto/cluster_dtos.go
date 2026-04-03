package dto

import "time"

// @MappedFrom AddMemberDTO(String nodeId, String zone, String name, NodeType nodeType, String version, boolean isSeed, boolean isLeaderEligible, Date registrationDate, int priority, String memberHost, int memberPort, String adminApiAddress, String wsAddress, String tcpAddress, String udpAddress, boolean isActive, boolean isHealthy)
type AddMemberDTO struct {
	NodeID           string    `json:"nodeId"`
	Zone             string    `json:"zone"`
	Name             string    `json:"name"`
	NodeType         string    `json:"nodeType"` // Or specific node type enum
	Version          string    `json:"version"`
	IsSeed           bool      `json:"isSeed"`
	IsLeaderEligible bool      `json:"isLeaderEligible"`
	RegistrationDate time.Time `json:"registrationDate"`
	Priority         int       `json:"priority"`
	MemberHost       string    `json:"memberHost"`
	MemberPort       int       `json:"memberPort"`
	AdminApiAddress  string    `json:"adminApiAddress"`
	WsAddress        string    `json:"wsAddress"`
	TcpAddress       string    `json:"tcpAddress"`
	UdpAddress       string    `json:"udpAddress"`
	IsActive         bool      `json:"isActive"`
	IsHealthy        bool      `json:"isHealthy"`
}

// @MappedFrom UpdateMemberDTO(String zone, String name, Boolean isSeed, Boolean isLeaderEligible, Boolean isActive, Integer priority)
type UpdateMemberDTO struct {
	Zone             *string `json:"zone"`
	Name             *string `json:"name"`
	IsSeed           *bool   `json:"isSeed"`
	IsLeaderEligible *bool   `json:"isLeaderEligible"`
	IsActive         *bool   `json:"isActive"`
	Priority         *int    `json:"priority"`
}

// @MappedFrom SettingsDTO(int schemaVersion, Map<String, Object> settings)
type SettingsDTO struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Settings      map[string]interface{} `json:"settings"`
}
