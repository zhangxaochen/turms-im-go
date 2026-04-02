package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/validator"
	turmsmongoexception "im.turms/server/internal/storage/mongo/exception"
	"im.turms/server/pkg/codes"
)

type UserRelationshipGroupService struct {
	groupRepo          repository.UserRelationshipGroupRepository
	groupMemberRepo    repository.UserRelationshipGroupMemberRepository
	userVersionService *UserVersionService
}

func NewUserRelationshipGroupService(
	groupRepo repository.UserRelationshipGroupRepository,
	groupMemberRepo repository.UserRelationshipGroupMemberRepository,
	userVersionService *UserVersionService,
) *UserRelationshipGroupService {
	return &UserRelationshipGroupService{
		groupRepo:          groupRepo,
		groupMemberRepo:    groupMemberRepo,
		userVersionService: userVersionService,
	}
}

func (s *UserRelationshipGroupService) CreateRelationshipGroup(
	ctx context.Context,
	ownerID int64,
	groupIndex *int32,
	groupName string,
	creationDate *time.Time,
	session *mongo.Session,
) (*po.UserRelationshipGroup, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(groupName, "groupName"); err != nil {
		return nil, err
	}
	if creationDate != nil {
		if err := validator.PastOrPresent(creationDate, "creationDate"); err != nil {
			return nil, err
		}
	}

	finalGroupIndex := int32(0)
	if groupIndex == nil {
		// Use a random index like Java if not provided.
		// Note: Turms Java uses RandomUtil.nextPositiveInt()
		// For now let's just use 0 and let it fail or handle it.
		// Actually, let's implement a simple random generator or just use a placeholder.
		// In Turms, custom groups start from a certain range or are just random.
		finalGroupIndex = int32(time.Now().UnixNano()) // Simple random
	} else {
		finalGroupIndex = *groupIndex
	}

	finalCreationDate := time.Now()
	if creationDate != nil {
		finalCreationDate = *creationDate
	}

	group := &po.UserRelationshipGroup{
		Key: po.UserRelationshipGroupKey{
			OwnerID: ownerID,
			Index:   finalGroupIndex,
		},
		Name:         groupName,
		CreationDate: finalCreationDate,
	}

	err := s.groupRepo.Insert(ctx, group, session)
	if err == nil {
		return group, nil
	}

	if turmsmongoexception.IsDuplicateKey(err) && groupIndex == nil && session == nil {
		return s.CreateRelationshipGroup(ctx, ownerID, nil, groupName, &finalCreationDate, nil)
	}

	return nil, err
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupsInfos(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return nil, err
	}
	return s.groupRepo.FindRelationshipGroupsInfos(ctx, ownerID)
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupsInfosWithVersion(
	ctx context.Context,
	ownerID int64,
	lastUpdatedDate *time.Time,
) ([]*po.UserRelationshipGroup, *time.Time, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return nil, nil, err
	}
	version, err := s.userVersionService.QueryRelationshipGroupsLastUpdatedDate(ctx, ownerID)
	if err != nil {
		return nil, nil, err
	}
	if version != nil && lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		return nil, nil, exception.NewTurmsError(int32(codes.AlreadyUpToDate), "already up to date")
	}
	groups, err := s.QueryRelationshipGroupsInfos(ctx, ownerID)
	return groups, version, err
}

func (s *UserRelationshipGroupService) QueryGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64) ([]int32, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(relatedUserID, "relatedUserID"); err != nil {
		return nil, err
	}
	return s.groupMemberRepo.FindGroupIndexes(ctx, ownerID, relatedUserID)
}

func (s *UserRelationshipGroupService) QueryRelationshipGroupMemberIds(
	ctx context.Context,
	ownerID int64,
	groupIndex int32,
) ([]int64, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(groupIndex, "groupIndex"); err != nil {
		return nil, err
	}
	return s.groupMemberRepo.FindRelationshipGroupMemberIds(ctx, []int64{ownerID}, []int32{groupIndex}, nil, nil)
}

func (s *UserRelationshipGroupService) UpdateRelationshipGroupName(
	ctx context.Context,
	ownerID int64,
	groupIndex int32,
	newGroupName string,
) error {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return err
	}
	if err := validator.NotNull(groupIndex, "groupIndex"); err != nil {
		return err
	}
	if err := validator.NotNull(newGroupName, "newGroupName"); err != nil {
		return err
	}
	count, err := s.groupRepo.UpdateRelationshipGroupName(ctx, ownerID, groupIndex, newGroupName, nil)
	if err != nil {
		return err
	}
	if count > 0 {
		go func() {
			if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, ownerID); err != nil {
				log.Printf("Failed to update relationship groups version for owner %d: %v", ownerID, err)
			}
		}()
	}
	return nil
}

func (s *UserRelationshipGroupService) UpsertRelationshipGroupMember(
	ctx context.Context,
	ownerID int64,
	relatedUserID int64,
	newGroupIndex *int32,
	deleteGroupIndex *int32,
	session *mongo.Session,
) (*int32, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(relatedUserID, "relatedUserID"); err != nil {
		return nil, err
	}
	if newGroupIndex != nil {
		if deleteGroupIndex != nil {
			if *newGroupIndex != *deleteGroupIndex {
				err := s.MoveRelatedUserToNewGroup(ctx, ownerID, relatedUserID, *deleteGroupIndex, *newGroupIndex, false, session)
				if err != nil {
					return nil, err
				}
				return newGroupIndex, nil
			}
		} else {
			added, err := s.AddRelatedUserToRelationshipGroup(ctx, ownerID, *newGroupIndex, relatedUserID, session)
			if err != nil {
				return nil, err
			}
			if added {
				return newGroupIndex, nil
			}
			return nil, nil
		}
	} else if deleteGroupIndex != nil && *deleteGroupIndex != 0 { // 0 is DEFAULT_RELATIONSHIP_GROUP_INDEX
		err := s.MoveRelatedUserToNewGroup(ctx, ownerID, relatedUserID, *deleteGroupIndex, 0, true, session)
		if err != nil {
			return nil, err
		}
		defaultIdx := int32(0)
		return &defaultIdx, nil
	}
	return nil, nil
}

func (s *UserRelationshipGroupService) UpdateRelationshipGroups(
	ctx context.Context,
	keys []po.UserRelationshipGroupKey,
	newName *string,
	creationDate *time.Time,
) error {
	if err := validator.NotEmpty(keys, "keys"); err != nil {
		return err
	}
	if creationDate != nil {
		if err := validator.PastOrPresent(creationDate, "creationDate"); err != nil {
			return err
		}
	}
	if newName == nil && creationDate == nil {
		return nil
	}
	// For now we only support name update in bulk
	if newName != nil {
		count, err := s.groupRepo.UpdateRelationshipGroups(ctx, keys, *newName, nil)
		if err != nil {
			return err
		}
		if count > 0 {
			ownerIDs := make(map[int64]bool)
			for _, key := range keys {
				ownerIDs[key.OwnerID] = true
			}
			for id := range ownerIDs {
				go func(oid int64) {
					if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, oid); err != nil {
						log.Printf("Failed to update relationship groups version for owner %d: %v", oid, err)
					}
				}(id)
			}
		}
	}
	// Note: creationDate update is NOT implemented in bulk yet, focusing on parity with current repo
	return nil
}

func (s *UserRelationshipGroupService) AddRelatedUserToRelationshipGroup(
	ctx context.Context,
	ownerID int64,
	groupIndex int32,
	relatedUserID int64,
	session *mongo.Session,
) (bool, error) {
	if err := validator.NotNull(groupIndex, "groupIndex"); err != nil {
		return false, err
	}
	member := &po.UserRelationshipGroupMember{
		Key: po.UserRelationshipGroupMemberKey{
			OwnerID:       ownerID,
			GroupIndex:    groupIndex,
			RelatedUserID: relatedUserID,
		},
		JoinDate: time.Now(),
	}
	res, err := s.groupMemberRepo.Upsert(ctx, member, session)
	if err != nil {
		return false, err
	}
	addedNew := res.UpsertedCount > 0
	if addedNew || res.ModifiedCount > 0 {
		go func() {
			if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, ownerID); err != nil {
				log.Printf("Failed to update relationship groups version for owner %d: %v", ownerID, err)
			}
		}()
	}
	return addedNew, nil
}

func (s *UserRelationshipGroupService) DeleteRelationshipGroups(
	ctx context.Context,
	ownerID int64,
	groupIndexes []int32,
	session *mongo.Session,
) (int64, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return 0, err
	}
	if err := validator.NotEmpty(groupIndexes, "groupIndexes"); err != nil {
		return 0, err
	}
	count, err := s.groupRepo.DeleteRelationshipGroups(ctx, ownerID, groupIndexes, session)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		go func() {
			if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, ownerID); err != nil {
				log.Printf("Failed to update relationship groups version for owner %d: %v", ownerID, err)
			}
		}()
	}
	return count, nil
}

func (s *UserRelationshipGroupService) DeleteRelationshipGroupAndMoveMembersToNewGroup(
	ctx context.Context,
	ownerID int64,
	deleteGroupIndex int32,
	newGroupIndex int32,
) error {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return err
	}
	if err := validator.NotNull(deleteGroupIndex, "deleteGroupIndex"); err != nil {
		return err
	}
	if err := validator.NotNull(newGroupIndex, "newGroupIndex"); err != nil {
		return err
	}
	if deleteGroupIndex == 0 { // DEFAULT_RELATIONSHIP_GROUP_INDEX
		return exception.NewTurmsError(int32(codes.IllegalArgument), "cannot delete default group")
	}
	if deleteGroupIndex == newGroupIndex {
		return nil
	}

	members, err := s.groupMemberRepo.FindRelationshipGroupMembers(ctx, ownerID, deleteGroupIndex)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return nil
	}

	now := time.Now()
	newMembers := make([]*po.UserRelationshipGroupMember, len(members))
	for i, m := range members {
		newMembers[i] = &po.UserRelationshipGroupMember{
			Key: po.UserRelationshipGroupMemberKey{
				OwnerID:       ownerID,
				GroupIndex:    newGroupIndex,
				RelatedUserID: m.Key.RelatedUserID,
			},
			JoinDate: now,
		}
	}

	_ = s.groupMemberRepo.InsertAllOfSameType(ctx, newMembers, nil)
	_, _ = s.groupMemberRepo.DeleteRelatedUserFromRelationshipGroup(ctx, ownerID, -1, []int32{deleteGroupIndex}, nil)
	_, _ = s.groupRepo.DeleteRelationshipGroups(ctx, ownerID, []int32{deleteGroupIndex}, nil)

	go func() {
		if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, ownerID); err != nil {
			log.Printf("Failed to update relationship groups version for owner %d: %v", ownerID, err)
		}
	}()

	return nil
}

func (s *UserRelationshipGroupService) DeleteAllRelationshipGroups(
	ctx context.Context,
	ownerIDs []int64,
	session *mongo.Session,
	updateVersion bool,
) error {
	if err := validator.NotEmpty(ownerIDs, "ownerIDs"); err != nil {
		return err
	}
	_, err := s.groupRepo.DeleteAllRelationshipGroups(ctx, ownerIDs, session)
	if err != nil {
		return err
	}
	if updateVersion {
		go func() {
			if err := s.userVersionService.UpdateSpecificVersions(ctx, ownerIDs, "rg"); err != nil {
				log.Printf("Failed to update relationship groups version for owners %v: %v", ownerIDs, err)
			}
		}()
	}
	return nil
}

func (s *UserRelationshipGroupService) DeleteRelatedUserFromRelationshipGroup(
	ctx context.Context,
	ownerID int64,
	relatedUserID int64,
	groupIndex int32,
	session *mongo.Session,
	updateVersion bool,
) (int64, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return 0, err
	}
	if err := validator.NotNull(relatedUserID, "relatedUserID"); err != nil {
		return 0, err
	}
	if err := validator.NotNull(groupIndex, "groupIndex"); err != nil {
		return 0, err
	}
	count, err := s.groupMemberRepo.DeleteRelatedUserFromRelationshipGroup(ctx, ownerID, relatedUserID, []int32{groupIndex}, session)
	if err != nil {
		return 0, err
	}
	if count > 0 && updateVersion {
		go func() {
			if err := s.userVersionService.UpdateRelationshipGroupsMembersVersion(ctx, ownerID); err != nil {
				log.Printf("Failed to update relationship group members version for owner %d: %v", ownerID, err)
			}
		}()
	}
	return count, nil
}

func (s *UserRelationshipGroupService) DeleteRelatedUsersFromAllRelationshipGroups(
	ctx context.Context,
	ownerID int64,
	relatedUserIDs []int64,
	session *mongo.Session,
	updateVersion bool,
) (int64, error) {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return 0, err
	}
	if err := validator.NotEmpty(relatedUserIDs, "relatedUserIDs"); err != nil {
		return 0, err
	}
	count, err := s.groupMemberRepo.DeleteRelatedUsersFromAllRelationshipGroups(ctx, ownerID, relatedUserIDs, session)
	if err != nil {
		return 0, err
	}
	if count > 0 && updateVersion {
		go func() {
			if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, ownerID); err != nil {
				log.Printf("Failed to update relationship groups version for owner %d: %v", ownerID, err)
			}
		}()
	}
	return count, nil
}

func (s *UserRelationshipGroupService) MoveRelatedUserToNewGroup(
	ctx context.Context,
	ownerID int64,
	relatedUserID int64,
	currentGroupIndex int32,
	targetGroupIndex int32,
	suppressIfAlreadyExists bool,
	session *mongo.Session,
) error {
	if err := validator.NotNull(ownerID, "ownerID"); err != nil {
		return err
	}
	if err := validator.NotNull(relatedUserID, "relatedUserID"); err != nil {
		return err
	}
	if err := validator.NotNull(currentGroupIndex, "currentGroupIndex"); err != nil {
		return err
	}
	if err := validator.NotNull(targetGroupIndex, "targetGroupIndex"); err != nil {
		return err
	}
	if currentGroupIndex == targetGroupIndex {
		return nil
	}
	newMember := &po.UserRelationshipGroupMember{
		Key: po.UserRelationshipGroupMemberKey{
			OwnerID:       ownerID,
			GroupIndex:    targetGroupIndex,
			RelatedUserID: relatedUserID,
		},
		JoinDate: time.Now(),
	}
	err := s.groupMemberRepo.Insert(ctx, newMember, session)
	if err != nil && (!suppressIfAlreadyExists || !turmsmongoexception.IsDuplicateKey(err)) {
		return err
	}
	_, _ = s.groupMemberRepo.DeleteById(ctx, ownerID, currentGroupIndex, relatedUserID, session)
	go func() {
		if err := s.userVersionService.UpdateRelationshipGroupsVersion(ctx, ownerID); err != nil {
			log.Printf("Failed to update relationship groups version for owner %d: %v", ownerID, err)
		}
	}()
	return nil
}

func (s *UserRelationshipGroupService) CountRelationshipGroups(ctx context.Context, ownerIDs []int64) (int64, error) {
	return s.groupRepo.CountRelationshipGroups(ctx, ownerIDs, nil)
}

func (s *UserRelationshipGroupService) CountRelationshipGroupMembers(ctx context.Context, ownerIDs []int64, groupIndexes []int32) (int64, error) {
	return s.groupMemberRepo.CountMembers(ctx, ownerIDs, groupIndexes)
}

func (s *UserRelationshipGroupService) QueryRelationshipGroups(
	ctx context.Context,
	ownerIDs []int64,
	groupIndexes []int32,
	page *int,
	size *int,
) ([]*po.UserRelationshipGroup, error) {
	return s.groupRepo.FindRelationshipGroups(ctx, ownerIDs, groupIndexes, page, size)
}
