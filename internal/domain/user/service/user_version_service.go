package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/user/repository"
)

type UserVersionService struct {
	versionRepo repository.UserVersionRepository
}

func NewUserVersionService(versionRepo repository.UserVersionRepository) *UserVersionService {
	return &UserVersionService{
		versionRepo: versionRepo,
	}
}

func (s *UserVersionService) UpsertEmptyUserVersion(ctx context.Context, userID int64) error {
	return s.versionRepo.UpsertEmptyUserVersion(ctx, userID)
}

// @MappedFrom queryRelationshipsLastUpdatedDate(@NotNull Long userId)
func (s *UserVersionService) QueryRelationshipsLastUpdatedDate(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.Relationships, nil
}

func (s *UserVersionService) QuerySentGroupInvitationsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.SentGroupInvitations, nil
}

func (s *UserVersionService) QueryReceivedGroupInvitationsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.ReceivedGroupInvitations, nil
}

// @MappedFrom queryGroupJoinRequestsVersion(@NotNull Long userId)
// @MappedFrom queryGroupJoinRequestsVersion(@NotNull Long groupId)
func (s *UserVersionService) QueryGroupJoinRequestsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.GroupJoinRequests, nil
}

// @MappedFrom queryRelationshipGroupsLastUpdatedDate(@NotNull Long userId)
func (s *UserVersionService) QueryRelationshipGroupsLastUpdatedDate(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.RelationshipGroups, nil
}

// @MappedFrom queryJoinedGroupVersion(@NotNull Long userId)
func (s *UserVersionService) QueryJoinedGroupVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.JoinedGroups, nil
}

// @MappedFrom querySentFriendRequestsVersion(@NotNull Long userId)
func (s *UserVersionService) QuerySentFriendRequestsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.SentFriendRequests, nil
}

// @MappedFrom queryReceivedFriendRequestsVersion(@NotNull Long userId)
func (s *UserVersionService) QueryReceivedFriendRequestsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := s.versionRepo.FindUserVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return &v.ReceivedFriendRequests, nil
}

// @MappedFrom updateRelationshipsVersion(@NotEmpty Set<Long> userIds, @Nullable ClientSession session)
// @MappedFrom updateRelationshipsVersion(@NotNull Long userId, @Nullable ClientSession session)
func (s *UserVersionService) UpdateRelationshipsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"r": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateSentFriendRequestsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateSentFriendRequestsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"sfr": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateReceivedFriendRequestsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateReceivedFriendRequestsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"rfr": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateRelationshipGroupsVersion(@NotEmpty Set<Long> userIds)
// @MappedFrom updateRelationshipGroupsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateRelationshipGroupsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"rg": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateRelationshipGroupsMembersVersion(@NotEmpty Set<Long> userIds)
// @MappedFrom updateRelationshipGroupsMembersVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateRelationshipGroupsMembersVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"rgm": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateSentGroupInvitationsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateSentGroupInvitationsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"sgi": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateReceivedGroupInvitationsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateReceivedGroupInvitationsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"rgi": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateSentGroupJoinRequestsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateSentGroupJoinRequestsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"gjr": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

// @MappedFrom updateJoinedGroupsVersion(@NotNull Long userId)
func (s *UserVersionService) UpdateJoinedGroupsVersion(ctx context.Context, userID int64) error {
	now := time.Now()
	update := map[string]interface{}{"$set": map[string]interface{}{"jg": now}}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

func (s *UserVersionService) UpdateSpecificVersion(ctx context.Context, userID int64, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	now := time.Now()
	setFields := make(map[string]interface{})
	for _, f := range fields {
		setFields[f] = now
	}
	update := map[string]interface{}{"$set": setFields}
	return s.versionRepo.UpdateUserVersion(ctx, userID, update)
}

func (s *UserVersionService) UpdateSpecificVersions(ctx context.Context, userIDs []int64, fields ...string) error {
	if len(fields) == 0 || len(userIDs) == 0 {
		return nil
	}
	now := time.Now()
	setFields := make(map[string]interface{})
	for _, f := range fields {
		setFields[f] = now
	}
	update := map[string]interface{}{"$set": setFields}
	return s.versionRepo.UpdateUserVersions(ctx, userIDs, update)
}

func (s *UserVersionService) Delete(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return s.versionRepo.DeleteUserVersions(ctx, userIDs)
}
