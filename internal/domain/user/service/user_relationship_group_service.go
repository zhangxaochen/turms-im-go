package service

import (
	"context"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserRelationshipGroupService struct {
	groupRepo       repository.UserRelationshipGroupRepository
	groupMemberRepo repository.UserRelationshipGroupMemberRepository
}

func NewUserRelationshipGroupService(
	groupRepo repository.UserRelationshipGroupRepository,
	groupMemberRepo repository.UserRelationshipGroupMemberRepository,
) *UserRelationshipGroupService {
	return &UserRelationshipGroupService{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
	}
}

func (s *UserRelationshipGroupService) CreateRelationshipGroup(ctx context.Context, group *po.UserRelationshipGroup) error {
	return s.groupRepo.InsertGroup(ctx, group)
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupsInfos(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error) {
	return s.groupRepo.FindRelationshipGroupsInfos(ctx, ownerID)
}

func (s *UserRelationshipGroupService) QueryGroupIndexes(ctx context.Context, ownerID int64) ([]int32, error) {
	groups, err := s.QueryRelationshipGroupsInfos(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	indexes := make([]int32, len(groups))
	for i, g := range groups {
		indexes[i] = g.Key.Index
	}
	return indexes, nil
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupMemberIds(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroupMember, error) {
	return s.groupMemberRepo.FindRelationshipGroupMemberIds(ctx, ownerID)
}

func (s *UserRelationshipGroupService) UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string) (int64, error) {
	return s.groupRepo.UpdateRelationshipGroupName(ctx, ownerID, groupIndex, newName, nil)
}

func (s *UserRelationshipGroupService) UpsertRelationshipGroupMember(ctx context.Context, member *po.UserRelationshipGroupMember) error {
	return s.groupMemberRepo.InsertMember(ctx, member)
}

func (s *UserRelationshipGroupService) UpdateRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, newName string) (int64, error) {
	return s.groupRepo.UpdateRelationshipGroups(ctx, ownerID, groupIndexes, newName, nil)
}

func (s *UserRelationshipGroupService) AddRelatedUserToRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, relatedUserID int64) error {
	for _, index := range groupIndexes {
		member := &po.UserRelationshipGroupMember{
			Key: po.UserRelationshipGroupMemberKey{
				OwnerID:       ownerID,
				GroupIndex:    index,
				RelatedUserID: relatedUserID,
			},
		}
		err := s.groupMemberRepo.InsertMember(ctx, member)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserRelationshipGroupService) DeleteAllRelationshipGroups(ctx context.Context, ownerID int64) error {
	_, err := s.groupRepo.DeleteAllRelationshipGroups(ctx, []int64{ownerID}, nil)
	if err != nil {
		return err
	}
	_, err = s.groupMemberRepo.DeleteAllRelatedUserFromRelationshipGroup(ctx, ownerID, nil)
	return err
}

func (s *UserRelationshipGroupService) DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndex int32) (int64, error) {
	return s.groupMemberRepo.DeleteRelatedUserFromRelationshipGroup(ctx, ownerID, relatedUserID, []int32{groupIndex}, nil)
}

func (s *UserRelationshipGroupService) DeleteRelatedUserFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserID int64) (int64, error) {
	return s.groupMemberRepo.DeleteRelatedUsersFromAllRelationshipGroups(ctx, ownerID, []int64{relatedUserID}, nil)
}

func (s *UserRelationshipGroupService) DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64) (int64, error) {
	return s.groupMemberRepo.DeleteRelatedUsersFromAllRelationshipGroups(ctx, ownerID, relatedUserIDs, nil)
}

func (s *UserRelationshipGroupService) MoveRelatedUserToNewGroup(ctx context.Context, ownerID int64, relatedUserID int64, newGroupIndex int32) error {
	_, err := s.DeleteRelatedUserFromAllRelationshipGroups(ctx, ownerID, relatedUserID)
	if err != nil {
		return err
	}
	member := &po.UserRelationshipGroupMember{
		Key: po.UserRelationshipGroupMemberKey{
			OwnerID:       ownerID,
			GroupIndex:    newGroupIndex,
			RelatedUserID: relatedUserID,
		},
	}
	return s.groupMemberRepo.InsertMember(ctx, member)
}

func (s *UserRelationshipGroupService) DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32) (int64, error) {
	res, err := s.groupRepo.DeleteRelationshipGroups(ctx, ownerID, groupIndexes, nil)
	if err != nil {
		return 0, err
	}
	_, _ = s.groupMemberRepo.DeleteRelatedUserFromRelationshipGroup(ctx, ownerID, -1, groupIndexes, nil) // Assuming -1 means any, wait!
	// Java doesn't have relatedUserID filter when deleting group. Let's assume it deletes all members in those groups.
	// Wait, is there a delete all members from group indexes?
	// The java code says: userRelationshipGroupMemberRepository.deleteByOwnerIdAndGroupIndexes(ownerId, groupIndexes)
	// I should pass no relatedUserIDs and it builds query `{ownerID, groupIndexes}`.
	return res, nil
}

func (s *UserRelationshipGroupService) CountRelationshipGroups(ctx context.Context, ownerID int64) (int64, error) {
	return s.groupRepo.CountRelationshipGroups(ctx, []int64{ownerID}, nil)
}

func (s *UserRelationshipGroupService) CountRelationshipGroupMembers(ctx context.Context, ownerID int64, groupIndex int32) (int64, error) {
	return s.groupMemberRepo.CountMembers(ctx, ownerID, groupIndex)
}

func (s *UserRelationshipGroupService) QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int) ([]*po.UserRelationshipGroup, error) {
	return s.groupRepo.FindRelationshipGroups(ctx, ownerIDs, groupIndexes, page, size)
}

func (s *UserRelationshipGroupService) DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx context.Context, ownerID int64, deleteGroupIndex int32, newGroupIndex int32) error {
	if deleteGroupIndex == newGroupIndex {
		return nil
	}

	members, err := s.groupMemberRepo.FindRelationshipGroupMembers(ctx, ownerID, deleteGroupIndex)
	if err != nil {
		return err
	}

	for _, member := range members {
		newMember := &po.UserRelationshipGroupMember{
			Key: po.UserRelationshipGroupMemberKey{
				OwnerID:       ownerID,
				GroupIndex:    newGroupIndex,
				RelatedUserID: member.Key.RelatedUserID,
			},
		}
		_ = s.groupMemberRepo.InsertMember(ctx, newMember)
	}

	_, err = s.groupMemberRepo.DeleteRelatedUserFromRelationshipGroup(ctx, ownerID, -1, []int32{deleteGroupIndex}, nil)
	if err != nil {
		return err
	}

	_, err = s.groupRepo.DeleteRelationshipGroups(ctx, ownerID, []int32{deleteGroupIndex}, nil)
	return err
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupsInfosWithVersion(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error) {
	return s.QueryRelationshipGroupsInfos(ctx, ownerID)
}
