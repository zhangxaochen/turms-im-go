package service

import (
	"context"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupTypeService struct {
	groupTypeRepo *repository.GroupTypeRepository
}

func NewGroupTypeService(groupTypeRepo *repository.GroupTypeRepository) *GroupTypeService {
	return &GroupTypeService{
		groupTypeRepo: groupTypeRepo,
	}
}

// FindGroupType retrieves a group type by its ID.
func (s *GroupTypeService) FindGroupType(ctx context.Context, typeID int64) (*po.GroupType, error) {
	return s.groupTypeRepo.FindGroupType(ctx, typeID)
}

// EnsureDefaultGroupType creates the default group type if it does not exist.
func (s *GroupTypeService) EnsureDefaultGroupType(ctx context.Context) error {
	// Simple implementation for ensuring default group type.
	t, err := s.groupTypeRepo.FindGroupType(ctx, 0)
	if err != nil {
		return err
	}
	if t == nil {
		defaultType := &po.GroupType{
			ID:             0,
			Name:           "DEFAULT",
			GroupSizeLimit: 500,
		}
		return s.groupTypeRepo.InsertGroupType(ctx, defaultType)
	}
	return nil
}
