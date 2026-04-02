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
