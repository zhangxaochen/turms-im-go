package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrNotGroupOwner = errors.New("not the group owner")
)

type GroupService struct {
	groupRepo          *repository.GroupRepository
	groupMemberService *GroupMemberService
}

func NewGroupService(groupRepo *repository.GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

func (s *GroupService) SetGroupMemberService(groupMemberService *GroupMemberService) {
	s.groupMemberService = groupMemberService
}

// CreateGroup creates a new group.
func (s *GroupService) CreateGroup(ctx context.Context, creatorID, groupID int64, name, intro *string, minimumScore *int32) (*po.Group, error) {
	now := time.Now()
	group := &po.Group{
		ID:           groupID,
		CreatorID:    &creatorID,
		OwnerID:      &creatorID,
		Name:         name,
		Intro:        intro,
		MinimumScore: minimumScore,
		CreationDate: &now,
	}

	err := s.groupRepo.InsertGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// DeleteGroup performs a soft deletion of the group.
// Only the owner can delete the group.
func (s *GroupService) DeleteGroup(ctx context.Context, requesterID, groupID int64) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return ErrGroupNotFound
	}

	if *ownerID != requesterID {
		return ErrNotGroupOwner
	}

	now := time.Now()
	update := bson.M{}
	update["dd"] = now

	return s.groupRepo.UpdateGroup(ctx, groupID, update)
}

func (s *GroupService) QueryGroupTypeIdIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*int64, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil || group.DeletionDate != nil || (group.IsActive != nil && !*group.IsActive) {
		return nil, nil
	}
	return group.TypeID, nil
}

func (s *GroupService) AuthAndTransferGroupOwnership(
	ctx context.Context,
	requesterID, groupID, successorID int64,
	quitAfterTransfer bool,
	session mongo.SessionContext,
) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return ErrGroupNotFound
	}
	if *ownerID != requesterID {
		return ErrNotGroupOwner
	}
	if requesterID == successorID {
		return nil
	}

	// Update owner in repository
	update := bson.M{"oid": successorID}
	err = s.groupRepo.UpdateGroup(ctx, groupID, update)
	if err != nil {
		return err
	}

	// Update roles in group member repository
	// Successor becomes Owner, Requester becomes Member (if quitAfterTransfer is false)
	err = s.groupMemberService.UpdateGroupMemberRole(ctx, groupID, successorID, protocol.GroupMemberRole_OWNER, session)
	if err != nil {
		return err
	}

	if quitAfterTransfer {
		return s.groupMemberService.DeleteGroupMember(ctx, groupID, requesterID, session, false)
	} else {
		return s.groupMemberService.UpdateGroupMemberRole(ctx, groupID, requesterID, protocol.GroupMemberRole_MEMBER, session)
	}
}
// AuthAndDeleteGroup deletes a group after authorization check.
func (s *GroupService) AuthAndDeleteGroup(ctx context.Context, requesterID int64, groupID int64) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist")
	}
	if *ownerID != requesterID {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_TO_DELETE_GROUP), "Only the owner can delete the group")
	}
	return s.groupRepo.DeleteGroup(ctx, groupID)
}
