package controller

import (
	"context"

	"im.turms/server/internal/domain/gateway/session"
)

// SessionController handles HTTP admin API requests related to sessions.
type SessionController struct {
	sessionService *session.SessionService
}

func NewSessionController(sessionService *session.SessionService) *SessionController {
	return &SessionController{
		sessionService: sessionService,
	}
}

// DeleteSessions deletes sessions based on the provided user IDs and/or IPs.
// @MappedFrom deleteSessions(@QueryParam(required = false) Set<Long> ids, @QueryParam(required = false) Set<String> ips)
func (c *SessionController) DeleteSessions(ctx context.Context, ids []int64, ips []string) (int, error) {
	// In the future, this should integrate with Gin/Fiber HTTP request handlers
	// For now, it's a stub demonstrating the mapping
	count := 0
	// ... Logic omitted for now until SessionService.CloseLocalSessions is fully implemented in Go.
	return count, nil
}
