package controller

import (
	"fmt"

	clusterdto "im.turms/server/internal/domain/cluster/access/admin/dto"
	commoncontroller "im.turms/server/internal/domain/common/access/admin/controller"
	"im.turms/server/internal/domain/common/infra/cluster/discovery"
	"im.turms/server/internal/infra/property"
)

// MemberController maps to MemberController Java.
// @MappedFrom MemberController
type MemberController struct {
	*commoncontroller.BaseController
	discoveryService *discovery.DiscoveryService
}

func NewMemberController(propertiesManager *property.TurmsPropertiesManager, discoveryService *discovery.DiscoveryService) *MemberController {
	return &MemberController{
		BaseController:   commoncontroller.NewBaseController(propertiesManager),
		discoveryService: discoveryService,
	}
}

// @MappedFrom queryMembers
func (c *MemberController) QueryMembers() []clusterdto.MemberDTO {
	members := c.discoveryService.GetMembers()
	leaderID := c.discoveryService.GetLeaderID()
	dtos := make([]clusterdto.MemberDTO, len(members))
	for i, m := range members {
		dtos[i] = clusterdto.MemberDTO{
			NodeID:           m.NodeID,
			Zone:             m.Zone,
			Name:             m.Name,
			NodeType:         m.NodeType,
			IsSeed:           m.IsSeed,
			IsLeaderEligible: m.IsLeaderEligible,
			Priority:         m.Priority,
			MemberHost:       m.MemberHost,
			MemberPort:       m.MemberPort,
			AdminAPIAddress:  m.AdminAPIAddress,
			WsAddress:        m.WsAddress,
			TcpAddress:       m.TcpAddress,
			UdpAddress:       m.UdpAddress,
			IsActive:         m.IsActive,
			IsHealthy:        m.IsHealthy,
			IsLeader:         m.NodeID == leaderID,
		}
	}
	return dtos
}

// @MappedFrom removeMembers
func (c *MemberController) RemoveMembers(ids []string) error {
	return c.discoveryService.UnregisterMembers(ids)
}

// @MappedFrom addMember
func (c *MemberController) AddMember(addMemberDTO clusterdto.AddMemberDTO) error {
	if addMemberDTO.NodeType != nil && *addMemberDTO.NodeType != discovery.NodeTypeService && addMemberDTO.IsLeaderEligible != nil && *addMemberDTO.IsLeaderEligible {
		return fmt.Errorf("only turms-service servers can be the leader") // NodeType 0 is SERVICE
	}

	isLeaderEligible := false
	if addMemberDTO.IsLeaderEligible != nil {
		isLeaderEligible = *addMemberDTO.IsLeaderEligible
	}

	member := &discovery.Member{
		ClusterID:        c.discoveryService.GetLocalNodeID(), // mock
		NodeID:           addMemberDTO.NodeID,
		Zone:             addMemberDTO.Zone,
		Name:             addMemberDTO.Name,
		IsSeed:           addMemberDTO.IsSeed,
		IsLeaderEligible: isLeaderEligible,
		Priority:         addMemberDTO.Priority,
		IsActive:         false,
		IsHealthy:        false,
	}

	if addMemberDTO.NodeType != nil {
		member.NodeType = *addMemberDTO.NodeType
	}
	if addMemberDTO.MemberHost != nil {
		member.MemberHost = *addMemberDTO.MemberHost
	}

	return c.discoveryService.RegisterMember(member)
}

// @MappedFrom updateMember
func (c *MemberController) UpdateMember(id string, updateMemberDTO clusterdto.UpdateMemberDTO) error {
	return c.discoveryService.UpdateMemberInfo(id, updateMemberDTO.Zone, updateMemberDTO.Name, updateMemberDTO.IsSeed, updateMemberDTO.IsLeaderEligible, updateMemberDTO.IsActive, updateMemberDTO.Priority)
}

// @MappedFrom queryLeader
func (c *MemberController) QueryLeader() (*clusterdto.MemberDTO, error) {
	leaderID := c.discoveryService.GetLeaderID()
	if leaderID == "" {
		return nil, fmt.Errorf("NO_CONTENT")
	}
	m := c.discoveryService.GetMember(leaderID)
	if m == nil {
		return nil, fmt.Errorf("NO_CONTENT")
	}
	dto := &clusterdto.MemberDTO{
		NodeID:           m.NodeID,
		Zone:             m.Zone,
		Name:             m.Name,
		NodeType:         m.NodeType,
		IsSeed:           m.IsSeed,
		IsLeaderEligible: m.IsLeaderEligible,
		Priority:         m.Priority,
		MemberHost:       m.MemberHost,
		MemberPort:       m.MemberPort,
		AdminAPIAddress:  m.AdminAPIAddress,
		WsAddress:        m.WsAddress,
		TcpAddress:       m.TcpAddress,
		UdpAddress:       m.UdpAddress,
		IsActive:         m.IsActive,
		IsHealthy:        m.IsHealthy,
		IsLeader:         true,
	}
	return dto, nil
}

// @MappedFrom electNewLeader
func (c *MemberController) ElectNewLeader(id *string) error {
	if id == nil {
		return c.discoveryService.ElectNewLeaderByPriority()
	}
	return c.discoveryService.ElectNewLeaderByNodeID(*id)
}

// SettingController maps to SettingController Java.
// @MappedFrom SettingController
type SettingController struct {
	*commoncontroller.BaseController
}

func NewSettingController(propertiesManager *property.TurmsPropertiesManager) *SettingController {
	return &SettingController{
		BaseController: commoncontroller.NewBaseController(propertiesManager),
	}
}

// @MappedFrom queryClusterSettings
func (c *SettingController) QueryClusterSettings(queryLocalSettings bool, onlyMutable bool) clusterdto.SettingsDTO {
	var props *property.TurmsProperties
	if queryLocalSettings {
		props = c.PropertiesManager.GetLocalProperties()
	} else {
		props = c.PropertiesManager.GetGlobalProperties()
	}
	_ = props // TODO: actually use for setting construction

	// properties to map logic
	settings := make(map[string]interface{})
	// Mock implementation for now, should use reflection or generated code in production
	// but follows the Java pattern of manual conversion if needed.

	return clusterdto.SettingsDTO{
		SchemaVersion: "1.0", // TurmsProperties.SCHEMA_VERSION
		Settings:      settings,
	}
}

// @MappedFrom updateClusterSettings
func (c *SettingController) UpdateClusterSettings(reset bool, updateLocalSettings bool, turmsProperties map[string]interface{}) error {
	if updateLocalSettings {
		return c.PropertiesManager.UpdateLocalProperties(reset, turmsProperties)
	}
	return c.PropertiesManager.UpdateGlobalProperties(reset, turmsProperties)
}

// @MappedFrom queryClusterConfigMetadata
func (c *SettingController) QueryClusterConfigMetadata(queryLocalSettings bool, onlyMutable bool, withValue bool) clusterdto.SettingsDTO {
	metadata := make(map[string]interface{})
	// logic to get metadata
	return clusterdto.SettingsDTO{
		SchemaVersion: "1.0",
		Settings:      metadata,
	}
}
