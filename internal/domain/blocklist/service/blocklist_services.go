package service

import (
	"im.turms/server/internal/domain/blocklist/bo"
)

type BlocklistService struct {
}

func (s *BlocklistService) BlockIpStrings(ips []string, blockDuration int64) error {
	return nil
}

func (s *BlocklistService) UnblockIpStrings(ips []string) error {
	return nil
}

func (s *BlocklistService) UnblockAllIps() error {
	return nil
}

func (s *BlocklistService) GetBlockedIpStrings(ips []string) []bo.BlockedClient[string] {
	return nil
}

func (s *BlocklistService) GetBlockedIps(page, size int) []bo.BlockedClient[string] {
	return nil
}

func (s *BlocklistService) BlockUserIds(userIds []int64, blockDuration int64) error {
	return nil
}

func (s *BlocklistService) UnblockUserIds(userIds []int64) error {
	return nil
}

func (s *BlocklistService) UnblockAllUserIds() error {
	return nil
}

func (s *BlocklistService) GetBlockedUsers(userIds []int64) []bo.BlockedClient[int64] {
	return nil
}

func (s *BlocklistService) GetBlockedUsersByPage(page, size int) []bo.BlockedClient[int64] {
	return nil
}
