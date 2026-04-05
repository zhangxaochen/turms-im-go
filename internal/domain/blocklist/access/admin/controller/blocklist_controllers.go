package controller

import (
	"time"

	"im.turms/server/internal/domain/blocklist/access/admin/dto"
	"im.turms/server/internal/domain/blocklist/service"
	commoncontroller "im.turms/server/internal/domain/common/access/admin/controller"
)

// PaginatedBlockedIpsDTO wraps blocked IPs with a total count for pagination.
type PaginatedBlockedIpsDTO struct {
	Total      int64               `json:"total"`
	BlockedIps []dto.BlockedIpDTO `json:"blockedIps"`
}

// PaginatedBlockedUsersDTO wraps blocked users with a total count for pagination.
type PaginatedBlockedUsersDTO struct {
	Total         int64                 `json:"total"`
	BlockedUsers []dto.BlockedUserDTO `json:"blockedUsers"`
}

// IpBlocklistController maps to IpBlocklistController.java
// @MappedFrom IpBlocklistController
type IpBlocklistController struct {
	*commoncontroller.BaseController
	blocklistService *service.BlocklistService
}

// @MappedFrom addBlockedIps(@RequestBody AddBlockedIpsDTO addBlockedIpsDTO)
func (c *IpBlocklistController) AddBlockedIps(addBlockedIpsDTO dto.AddBlockedIpsDTO) error {
	return c.blocklistService.BlockIpStrings(addBlockedIpsDTO.IDs, addBlockedIpsDTO.BlockDurationMillis)
}

// @MappedFrom queryBlockedIps(Set<String> ids)
func (c *IpBlocklistController) QueryBlockedIpsByIds(ids []string) []dto.BlockedIpDTO {
	blockedClients := c.blocklistService.GetBlockedIpStrings(ids)
	dtos := make([]dto.BlockedIpDTO, len(blockedClients))
	for i, client := range blockedClients {
		dtos[i] = dto.BlockedIpDTO{
			ID:           client.ID,
			BlockEndTime: time.UnixMilli(client.BlockEndTimeMillis),
		}
	}
	return dtos
}

// @MappedFrom queryBlockedIps(int page, @QueryParam(required = false) Integer size)
// Bug fix: Added total count via CountBlockedIps() to match Java's paginated response
func (c *IpBlocklistController) QueryBlockedIpsByPage(page int, size *int) PaginatedBlockedIpsDTO {
	actualSize := c.GetPageSize(size)
	total, _ := c.blocklistService.CountBlockedIps()
	blockedClients := c.blocklistService.GetBlockedIps(page, actualSize)
	dtos := make([]dto.BlockedIpDTO, len(blockedClients))
	for i, client := range blockedClients {
		dtos[i] = dto.BlockedIpDTO{
			ID:           client.ID,
			BlockEndTime: time.UnixMilli(client.BlockEndTimeMillis),
		}
	}
	return PaginatedBlockedIpsDTO{
		Total:      total,
		BlockedIps: dtos,
	}
}

// @MappedFrom deleteBlockedIps(@QueryParam(required = false) Set<String> ids, @QueryParam(required = false) Boolean deleteAll)
// Bug fix: Added empty-ids guard — skip unblock when ids is empty and deleteAll is false
func (c *IpBlocklistController) DeleteBlockedIps(ids []string, deleteAll bool) error {
	if deleteAll {
		return c.blocklistService.UnblockAllIps()
	}
	if len(ids) == 0 {
		return nil
	}
	return c.blocklistService.UnblockIpStrings(ids)
}

// UserBlocklistController maps to UserBlocklistController.java
// @MappedFrom UserBlocklistController
type UserBlocklistController struct {
	*commoncontroller.BaseController
	blocklistService *service.BlocklistService
}

// @MappedFrom addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO)
func (c *UserBlocklistController) AddBlockedUserIds(addBlockedUserIdsDTO dto.AddBlockedUserIdsDTO) error {
	return c.blocklistService.BlockUserIds(addBlockedUserIdsDTO.IDs, addBlockedUserIdsDTO.BlockDurationMillis)
}

// @MappedFrom queryBlockedUserIds(Set<Long> ids)
func (c *UserBlocklistController) QueryBlockedUserIdsByIds(ids []int64) []dto.BlockedUserDTO {
	blockedClients := c.blocklistService.GetBlockedUsers(ids)
	dtos := make([]dto.BlockedUserDTO, len(blockedClients))
	for i, client := range blockedClients {
		dtos[i] = dto.BlockedUserDTO{
			ID:           client.ID,
			BlockEndTime: time.UnixMilli(client.BlockEndTimeMillis),
		}
	}
	return dtos
}

// @MappedFrom queryBlockedUserIds(int page, @QueryParam(required = false) Integer size)
// Bug fix: Added total count via CountBlockedUsers() to match Java's paginated response
func (c *UserBlocklistController) QueryBlockedUserIdsByPage(page int, size *int) PaginatedBlockedUsersDTO {
	actualSize := c.GetPageSize(size)
	total, _ := c.blocklistService.CountBlockedUsers()
	blockedClients := c.blocklistService.GetBlockedUsersByPage(page, actualSize)
	dtos := make([]dto.BlockedUserDTO, len(blockedClients))
	for i, client := range blockedClients {
		dtos[i] = dto.BlockedUserDTO{
			ID:           client.ID,
			BlockEndTime: time.UnixMilli(client.BlockEndTimeMillis),
		}
	}
	return PaginatedBlockedUsersDTO{
		Total:         total,
		BlockedUsers: dtos,
	}
}

// @MappedFrom deleteBlockedUserIds(@QueryParam(required = false) Set<Long> ids, @QueryParam(required = false) Boolean deleteAll)
// Bug fix: Added empty-ids guard — skip unblock when ids is empty and deleteAll is false
func (c *UserBlocklistController) DeleteBlockedUserIds(ids []int64, deleteAll bool) error {
	if deleteAll {
		return c.blocklistService.UnblockAllUserIds()
	}
	if len(ids) == 0 {
		return nil
	}
	return c.blocklistService.UnblockUserIds(ids)
}
