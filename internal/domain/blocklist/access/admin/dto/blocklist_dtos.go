package dto

import "time"

// @MappedFrom AddBlockedIpsDTO(Set<String> ids, long blockDurationMillis)
type AddBlockedIpsDTO struct {
	IDs                 []string `json:"ids"`
	BlockDurationMillis int64    `json:"blockDurationMillis"`
}

// @MappedFrom AddBlockedUserIdsDTO(Set<Long> ids, long blockDurationMillis)
type AddBlockedUserIdsDTO struct {
	IDs                 []int64 `json:"ids"`
	BlockDurationMillis int64   `json:"blockDurationMillis"`
}

// @MappedFrom BlockedIpDTO(String id, Date blockEndTime)
type BlockedIpDTO struct {
	ID           string    `json:"id"`
	BlockEndTime time.Time `json:"blockEndTime"`
}

// @MappedFrom BlockedUserDTO(Long id, Date blockEndTime)
type BlockedUserDTO struct {
	ID           int64     `json:"id"`
	BlockEndTime time.Time `json:"blockEndTime"`
}
