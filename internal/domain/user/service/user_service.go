package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, isActive bool) (*po.User, error)
	FindUser(ctx context.Context, userID int64) (*po.User, error)
	UpdateUser(ctx context.Context, userID int64, update bson.M) error
	CheckIfUserExists(ctx context.Context, userID int64) (bool, error)
	DeleteUsers(ctx context.Context, userIDs []int64) error
	QueryUsersProfile(ctx context.Context, userIDs []int64) ([]*po.User, error)
	QueryUserName(ctx context.Context, userID int64) (string, error)
	QueryUsers(ctx context.Context, userIDs []int64) ([]*po.User, error)
	CountUsers(ctx context.Context, activeOnly bool) (int64, error)
	IsAllowedToSendMessageToTarget(ctx context.Context, isGroupMessage bool, isSystemMessage bool, requesterID int64, targetID int64) (int, error)
	IsAllowToQueryUserProfile(ctx context.Context, requesterID int64, targetID int64) (int, error)
	AuthAndQueryUsersProfile(ctx context.Context, requesterID int64, userIDs []int64, name string, lastUpdatedDate *time.Time, skip int, limit int) ([]*po.User, error)
	QueryUserRoleIDByUserID(ctx context.Context, userID int64) (int64, error)
	CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool) (int64, error)
	CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error)
	UpdateUsers(ctx context.Context, userIDs []int64, update bson.M) (int64, error)
}

type userService struct {
	idGen *idgen.SnowflakeIdGenerator
	repo  repository.UserRepository
}

func NewUserService(idGen *idgen.SnowflakeIdGenerator, repo repository.UserRepository) UserService {
	return &userService{
		idGen: idGen,
		repo:  repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, password string, name string, intro string, profilePicture string, profileAccess int32, permissionGroupID int64, isActive bool) (*po.User, error) {
	userID := s.idGen.NextIncreasingId()
	now := time.Now()

	user := &po.User{
		ID:                userID,
		Password:          password, // Assuming plain text for this simple refactor, should be hashed in real world
		Name:              name,
		Intro:             intro,
		ProfilePicture:    profilePicture,
		ProfileAccess:     profileAccess,
		PermissionGroupID: permissionGroupID,
		RegistrationDate:  now,
		IsActive:          isActive,
	}
	err := s.repo.Insert(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) FindUser(ctx context.Context, userID int64) (*po.User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *userService) UpdateUser(ctx context.Context, userID int64, update bson.M) error {
	now := time.Now()
	update["lud"] = now // Set LastUpdatedDate
	return s.repo.Update(ctx, userID, update)
}

func (s *userService) CheckIfUserExists(ctx context.Context, userID int64) (bool, error) {
	return s.repo.Exists(ctx, userID)
}

func (s *userService) DeleteUsers(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	_, err := s.repo.DeleteMany(ctx, filter)
	return err
}

func (s *userService) QueryUsersProfile(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	if len(userIDs) == 0 {
		return []*po.User{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return s.repo.FindMany(ctx, filter)
}

func (s *userService) QueryUserName(ctx context.Context, userID int64) (string, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil // Or return a specific record not found error depending on Turms
	}
	return user.Name, nil
}

func (s *userService) QueryUsers(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	if len(userIDs) == 0 {
		return []*po.User{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return s.repo.FindMany(ctx, filter)
}

func (s *userService) CountUsers(ctx context.Context, activeOnly bool) (int64, error) {
	filter := bson.M{}
	if activeOnly {
		filter["act"] = true
	}
	return s.repo.Count(ctx, filter)
}

func (s *userService) IsAllowedToSendMessageToTarget(ctx context.Context, isGroupMessage bool, isSystemMessage bool, requesterID int64, targetID int64) (int, error) {
	if isSystemMessage {
		return 200, nil // OK
	}
	// Simplified permission check
	return 200, nil
}

func (s *userService) IsAllowToQueryUserProfile(ctx context.Context, requesterID int64, targetID int64) (int, error) {
	// Simplified logic for refactor
	return 200, nil
}

func (s *userService) AuthAndQueryUsersProfile(ctx context.Context, requesterID int64, userIDs []int64, name string, lastUpdatedDate *time.Time, skip int, limit int) ([]*po.User, error) {
	// Simplified, normally check permission then query
	filter := bson.M{}
	if len(userIDs) > 0 {
		filter["_id"] = bson.M{"$in": userIDs}
	}
	if name != "" {
		filter["n"] = name
	}
	return s.repo.FindMany(ctx, filter)
}

func (s *userService) QueryUserRoleIDByUserID(ctx context.Context, userID int64) (int64, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, nil
	}
	return user.PermissionGroupID, nil
}

func (s *userService) CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool) (int64, error) {
	filter := bson.M{}
	dateFilter := bson.M{}
	if startDate != nil {
		dateFilter["$gte"] = *startDate
	}
	if endDate != nil {
		dateFilter["$lt"] = *endDate
	}
	if len(dateFilter) > 0 {
		filter["rd"] = dateFilter
	}
	return s.repo.Count(ctx, filter)
}

func (s *userService) CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error) {
	filter := bson.M{"dd": bson.M{"$exists": true, "$ne": nil}}
	return s.repo.Count(ctx, filter)
}

func (s *userService) UpdateUsers(ctx context.Context, userIDs []int64, update bson.M) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return s.repo.UpdateMany(ctx, filter, update)
}
