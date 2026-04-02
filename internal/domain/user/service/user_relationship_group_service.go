package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
	filter := bson.M{"_id.oid": ownerID}
	return s.groupRepo.FindGroups(ctx, filter)
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

func (s *UserRelationshipGroupService) QueryRelationshipGroupMemberIds(ctx context.Context, filter bson.M) ([]*po.UserRelationshipGroupMember, error) {
	return s.groupMemberRepo.FindMembers(ctx, filter)
}

func (s *UserRelationshipGroupService) UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string) (int64, error) {
	filter := bson.M{"_id.oid": ownerID, "_id.i": groupIndex}
	update := bson.M{"$set": bson.M{"n": newName}}
	return s.groupRepo.UpdateGroups(ctx, filter, update)
}

func (s *UserRelationshipGroupService) UpsertRelationshipGroupMember(ctx context.Context, member *po.UserRelationshipGroupMember) error {
	return s.groupMemberRepo.UpsertMember(ctx, member)
}

func (s *UserRelationshipGroupService) UpdateRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, update bson.M) (int64, error) {
	filter := bson.M{"_id.oid": ownerID}
	if len(groupIndexes) > 0 {
		filter["_id.i"] = bson.M{"$in": groupIndexes}
	}
	return s.groupRepo.UpdateGroups(ctx, filter, update)
}

func (s *UserRelationshipGroupService) AddRelatedUserToRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, relatedUserID int64) error {
	// Loop to upsert
	for _, index := range groupIndexes {
		member := &po.UserRelationshipGroupMember{
			Key: po.UserRelationshipGroupMemberKey{
				OwnerID:       ownerID,
				GroupIndex:    index,
				RelatedUserID: relatedUserID,
			},
		}
		err := s.groupMemberRepo.UpsertMember(ctx, member)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserRelationshipGroupService) DeleteAllRelationshipGroups(ctx context.Context, ownerID int64) error {
	_, err := s.groupRepo.DeleteGroups(ctx, bson.M{"_id.oid": ownerID})
	if err != nil {
		return err
	}
	_, err = s.groupMemberRepo.DeleteMembers(ctx, bson.M{"_id.oid": ownerID})
	return err
}

func (s *UserRelationshipGroupService) DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndex int32) (int64, error) {
	filter := bson.M{"_id.oid": ownerID, "_id.rid": relatedUserID, "_id.gi": groupIndex}
	return s.groupMemberRepo.DeleteMembers(ctx, filter)
}

func (s *UserRelationshipGroupService) DeleteRelatedUserFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserID int64) (int64, error) {
	filter := bson.M{"_id.oid": ownerID, "_id.rid": relatedUserID}
	return s.groupMemberRepo.DeleteMembers(ctx, filter)
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
	return s.groupMemberRepo.UpsertMember(ctx, member)
}

func (s *UserRelationshipGroupService) DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32) (int64, error) {
	filter := bson.M{"_id.oid": ownerID, "_id.i": bson.M{"$in": groupIndexes}}
	res, err := s.groupRepo.DeleteGroups(ctx, filter)
	if err != nil {
		return 0, err
	}
	// Also delete members
	filterMembers := bson.M{"_id.oid": ownerID, "_id.gi": bson.M{"$in": groupIndexes}}
	_, _ = s.groupMemberRepo.DeleteMembers(ctx, filterMembers)
	return res, nil
}

func (s *UserRelationshipGroupService) CountRelationshipGroups(ctx context.Context, ownerID int64) (int64, error) {
	filter := bson.M{"_id.oid": ownerID}
	return s.groupRepo.CountGroups(ctx, filter)
}

func (s *UserRelationshipGroupService) CountRelationshipGroupMembers(ctx context.Context, filter bson.M) (int64, error) {
	return s.groupMemberRepo.CountMembers(ctx, filter)
}

func (s *UserRelationshipGroupService) DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64) (int64, error) {
	if len(relatedUserIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{"_id.oid": ownerID, "_id.rid": bson.M{"$in": relatedUserIDs}}
	return s.groupMemberRepo.DeleteMembers(ctx, filter)
}

func (s *UserRelationshipGroupService) QueryRelationshipGroups(ctx context.Context, filter bson.M) ([]*po.UserRelationshipGroup, error) {
	return s.groupRepo.FindGroups(ctx, filter)
}

func (s *UserRelationshipGroupService) DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx context.Context, ownerID int64, deleteGroupIndex int32, newGroupIndex int32) error {
	if deleteGroupIndex == newGroupIndex {
		return nil
	}

	filter := bson.M{"_id.oid": ownerID, "_id.gi": deleteGroupIndex}
	members, err := s.groupMemberRepo.FindMembers(ctx, filter)
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
		_ = s.groupMemberRepo.UpsertMember(ctx, newMember)
	}

	_, err = s.groupMemberRepo.DeleteMembers(ctx, filter)
	if err != nil {
		return err
	}

	filterGroup := bson.M{"_id.oid": ownerID, "_id.i": deleteGroupIndex}
	_, err = s.groupRepo.DeleteGroups(ctx, filterGroup)
	return err
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupsInfosWithVersion(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error) {
	// Note: Versioning integration to be added when UserVersionService is injected
	return s.QueryRelationshipGroupsInfos(ctx, ownerID)
}
