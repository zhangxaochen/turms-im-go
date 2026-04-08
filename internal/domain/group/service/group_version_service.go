package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/group/repository"
)

type GroupVersionService struct {
	groupVersionRepo *repository.GroupVersionRepository
}

func NewGroupVersionService(groupVersionRepo *repository.GroupVersionRepository) *GroupVersionService {
	return &GroupVersionService{
		groupVersionRepo: groupVersionRepo,
	}
}

func (s *GroupVersionService) InitVersions(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.InsertVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateMembersVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateMembersVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateInformationVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateInformationVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateBlocklistVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateBlocklistVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateInvitationsVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateInvitationsVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateJoinRequestsVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateJoinRequestsVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateJoinQuestionsVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateJoinQuestionsVersion(ctx, groupID)
}

// UpdateVersionFields conditionally updates multiple version fields in a single atomic operation.
// @MappedFrom updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)
func (s *GroupVersionService) UpdateVersionFields(ctx context.Context, groupID int64, updateMembers, updateBlocklist, updateJoinRequests, updateJoinQuestions bool) error {
	return s.groupVersionRepo.UpdateVersionFields(ctx, groupID, updateMembers, updateBlocklist, updateJoinRequests, updateJoinQuestions)
}

// UpdateMembersVersionForGroups batch-updates the members version for a set of group IDs.
// @MappedFrom updateMembersVersion(@Nullable Set<Long> groupIds)
func (s *GroupVersionService) UpdateMembersVersionForGroups(ctx context.Context, groupIDs []int64) error {
	return s.groupVersionRepo.UpdateVersions(ctx, groupIDs, "mbr")
}

// UpdateMembersVersionForAll updates the members version for ALL groups.
// @MappedFrom updateMembersVersion()
func (s *GroupVersionService) UpdateMembersVersionForAll(ctx context.Context) error {
	return s.groupVersionRepo.UpdateVersions(ctx, nil, "mbr")
}

// UpdateSpecificVersionForGroups updates a specific version field for a set of group IDs.
// @MappedFrom updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field)
func (s *GroupVersionService) UpdateSpecificVersionForGroups(ctx context.Context, groupIDs []int64, field string) error {
	return s.groupVersionRepo.UpdateVersions(ctx, groupIDs, field)
}

// UpdateSpecificVersionForAll updates a specific version field for ALL groups.
// @MappedFrom updateSpecificVersion(@NotNull String field)
func (s *GroupVersionService) UpdateSpecificVersionForAll(ctx context.Context, field string) error {
	return s.groupVersionRepo.UpdateVersions(ctx, nil, field)
}

// @MappedFrom queryGroupInvitationsVersion(@NotNull Long groupId)
func (s *GroupVersionService) QueryGroupInvitationsVersion(ctx context.Context, groupID int64) (*time.Time, error) {
	return s.groupVersionRepo.FindInvitations(ctx, groupID)
}

// @MappedFrom queryGroupJoinRequestsVersion(@NotNull Long groupId)
func (s *GroupVersionService) QueryGroupJoinRequestsVersion(ctx context.Context, groupID int64) (*time.Time, error) {
	return s.groupVersionRepo.FindJoinRequests(ctx, groupID)
}

// @MappedFrom queryGroupJoinQuestionsVersion(@NotNull Long groupId)
func (s *GroupVersionService) QueryGroupJoinQuestionsVersion(ctx context.Context, groupID int64) (*time.Time, error) {
	return s.groupVersionRepo.FindJoinQuestions(ctx, groupID)
}

// @MappedFrom queryMembersVersion(@NotNull Long groupId)
func (s *GroupVersionService) QueryGroupMembersVersion(ctx context.Context, groupID int64) (*time.Time, error) {
	return s.groupVersionRepo.FindMembers(ctx, groupID)
}

// @MappedFrom queryBlocklistVersion(@NotNull Long groupId)
func (s *GroupVersionService) QueryGroupBlocklistVersion(ctx context.Context, groupID int64) (*time.Time, error) {
	return s.groupVersionRepo.FindBlocklist(ctx, groupID)
}

// Upsert creates or updates all group version records.
func (s *GroupVersionService) Upsert(ctx context.Context, groupID int64, timestamp time.Time) error {
	return s.groupVersionRepo.Upsert(ctx, groupID, timestamp)
}

// Delete deletes group versions by group IDs.
func (s *GroupVersionService) Delete(ctx context.Context, groupIDs []int64) error {
	return s.groupVersionRepo.DeleteByIds(ctx, groupIDs)
}
