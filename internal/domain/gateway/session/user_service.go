package session

import (
	"context"
)

// UserService maps to UserService in Java for gateway session handling.
// @MappedFrom UserService
type UserService struct {
}

// @MappedFrom authenticate(@NotNull Long userId, @Nullable String rawPassword)
func (s *UserService) Authenticate(ctx context.Context, userID int64, rawPassword string) (bool, error) {
	// Stub implementation
	return false, nil
}
