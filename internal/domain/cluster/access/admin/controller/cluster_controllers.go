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
// Bug fix: Return deleted count to match Java's DeleteResultDTO
func (c *MemberController) RemoveMembers(ids []string) (int64, error) {
	err := c.discoveryService.UnregisterMembers(ids)
	if err != nil {
		return 0, err
	}
	return int64(len(ids)), nil
}

// @MappedFrom addMember
func (c *MemberController) AddMember(addMemberDTO clusterdto.AddMemberDTO) error {
	if addMemberDTO.NodeType != nil && *addMemberDTO.NodeType != discovery.NodeTypeService && addMemberDTO.IsLeaderEligible != nil && *addMemberDTO.IsLeaderEligible {
		return fmt.Errorf("only turms-service servers can be the leader")
	}

	isLeaderEligible := false
	if addMemberDTO.IsLeaderEligible != nil {
		isLeaderEligible = *addMemberDTO.IsLeaderEligible
	}

	member := &discovery.Member{
		// Bug fix: Use GetLocalClusterID() instead of GetLocalNodeID() for ClusterID
		ClusterID:        c.discoveryService.GetLocalClusterID(),
		NodeID:           addMemberDTO.NodeID,
		Zone:             addMemberDTO.Zone,
		Name:             addMemberDTO.Name,
		IsSeed:           addMemberDTO.IsSeed,
		IsLeaderEligible: isLeaderEligible,
		Priority:         addMemberDTO.Priority,
	}

	// Bug fix: Use DTO values for IsActive/IsHealthy instead of hardcoding false
	if addMemberDTO.IsActive != nil {
		member.IsActive = *addMemberDTO.IsActive
	}
	if addMemberDTO.IsHealthy != nil {
		member.IsHealthy = *addMemberDTO.IsHealthy
	}

	if addMemberDTO.NodeType != nil {
		member.NodeType = *addMemberDTO.NodeType
	}
	if addMemberDTO.MemberHost != nil {
		member.MemberHost = *addMemberDTO.MemberHost
	}

	// Bug fix: Populate missing fields from DTO
	if addMemberDTO.MemberPort != nil {
		member.MemberPort = *addMemberDTO.MemberPort
	}
	if addMemberDTO.AdminAPIAddress != nil {
		member.AdminAPIAddress = *addMemberDTO.AdminAPIAddress
	}
	if addMemberDTO.WsAddress != nil {
		member.WsAddress = *addMemberDTO.WsAddress
	}
	if addMemberDTO.TcpAddress != nil {
		member.TcpAddress = *addMemberDTO.TcpAddress
	}
	if addMemberDTO.UdpAddress != nil {
		member.UdpAddress = *addMemberDTO.UdpAddress
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
		return nil, nil
	}
	m := c.discoveryService.GetMember(leaderID)
	if m == nil {
		return nil, nil
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
// Bug fix: Return the elected leader member to match Java's return type
func (c *MemberController) ElectNewLeader(id *string) (*clusterdto.MemberDTO, error) {
	if id == nil {
		err := c.discoveryService.ElectNewLeaderByPriority()
		if err != nil {
			return nil, err
		}
	} else {
		err := c.discoveryService.ElectNewLeaderByNodeID(*id)
		if err != nil {
			return nil, err
		}
	}
	// Return the new leader
	leaderID := c.discoveryService.GetLeaderID()
	if leaderID == "" {
		return nil, nil
	}
	m := c.discoveryService.GetMember(leaderID)
	if m == nil {
		return nil, nil
	}
	return &clusterdto.MemberDTO{
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
	}, nil
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

// TurmsPropertiesSchemaVersion is the Go equivalent of TurmsProperties.SCHEMA_VERSION in Java.
const TurmsPropertiesSchemaVersion = 1

// @MappedFrom queryClusterSettings
func (c *SettingController) QueryClusterSettings(queryLocalSettings bool, onlyMutable bool) clusterdto.SettingsDTO {
	var props *property.TurmsProperties
	if queryLocalSettings {
		props = c.PropertiesManager.GetLocalProperties()
	} else {
		props = c.PropertiesManager.GetGlobalProperties()
	}

	// Bug fix: Use props for setting construction instead of discarding it
	settings := convertPropertiesToValueMap(props, onlyMutable)

	return clusterdto.SettingsDTO{
		SchemaVersion: TurmsPropertiesSchemaVersion,
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
	var props *property.TurmsProperties
	if queryLocalSettings {
		props = c.PropertiesManager.GetLocalProperties()
	} else {
		props = c.PropertiesManager.GetGlobalProperties()
	}

	// Get metadata (TODO: implement proper metadata lookup from TurmsPropertiesInspector)
	metadata := make(map[string]interface{})

	// Bug fix: Merge property values with metadata when withValue is true
	if withValue {
		valueMap := convertPropertiesToValueMap(props, onlyMutable)
		for k, v := range valueMap {
			metadata[k] = v
		}
	}

	return clusterdto.SettingsDTO{
		SchemaVersion: TurmsPropertiesSchemaVersion,
		Settings:      metadata,
	}
}

// convertPropertiesToValueMap converts TurmsProperties to a map of settings values.
// When onlyMutable is true, only mutable properties should be included.
// TODO: Implement full reflection-based conversion to match Java's TurmsPropertiesInspector
func convertPropertiesToValueMap(props *property.TurmsProperties, onlyMutable bool) map[string]interface{} {
	if props == nil {
		return make(map[string]interface{})
	}
	settings := make(map[string]interface{})
	// Flatten the nested property structs into a dot-separated key map.
	// This is a partial implementation; full conversion requires reflection or generated code.
	settings["service.adminApi.maxDayDifferencePerRequest"] = props.Service.AdminApi.MaxDayDifferencePerRequest
	settings["service.adminApi.maxHourDifferencePerCountRequest"] = props.Service.AdminApi.MaxHourDifferencePerCountRequest
	settings["service.adminApi.maxDayDifferencePerCountRequest"] = props.Service.AdminApi.MaxDayDifferencePerCountRequest
	settings["service.adminApi.maxMonthDifferencePerCountRequest"] = props.Service.AdminApi.MaxMonthDifferencePerCountRequest
	settings["service.adminApi.maxAvailableRecordsPerRequest"] = props.Service.AdminApi.MaxAvailableRecordsPerRequest
	settings["service.adminApi.defaultAvailableRecordsPerRequest"] = props.Service.AdminApi.DefaultAvailableRecordsPerRequest
	return settings
}
