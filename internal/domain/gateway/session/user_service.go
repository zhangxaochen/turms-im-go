package session

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"im.turms/server/internal/domain/user/repository"
	"im.turms/server/internal/pkg/security"
)

// UserService maps to UserService in Java for gateway session handling.
// @MappedFrom UserService
type UserService struct {
	userRepository repository.UserRepository
	enabled        bool // From properties or repo in Java
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	// Hardcoded enabled = true for now, in Java it checks userRepository.isEnabled()
	return &UserService{
		userRepository: userRepository,
		enabled:        true,
	}
}

// @MappedFrom authenticate(@NotNull Long userId, @Nullable String rawPassword)
func (s *UserService) Authenticate(ctx context.Context, userID int64, rawPassword string) (bool, error) {
	passwordHash, err := s.userRepository.FindPassword(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // user not found
		}
		return false, err
	}
	if passwordHash == nil {
		// Should not happen with our explicit FindPassword implementation, but for safety
		return rawPassword == "", nil
	}
	return security.MatchesPassword(rawPassword, *passwordHash), nil
}

// @MappedFrom isActiveAndNotDeleted(@NotNull Long userId)
func (s *UserService) IsActiveAndNotDeleted(ctx context.Context, userID int64) (bool, error) {
	return s.userRepository.IsActiveAndNotDeleted(ctx, userID)
}
