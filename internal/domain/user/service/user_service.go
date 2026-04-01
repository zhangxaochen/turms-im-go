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
