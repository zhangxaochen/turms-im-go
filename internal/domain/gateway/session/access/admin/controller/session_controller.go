package controller

import (
	"context"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/common/constant"
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
	count := 0
	if len(ids) > 0 {
		err := c.sessionService.CloseLocalSessionsByUserIds(ctx, ids, constant.SessionCloseStatus_DISCONNECTED_BY_ADMIN)
		if err != nil {
			return 0, err
		}
		count += len(ids)
	}

	if len(ips) > 0 {
		var byteIps [][]byte
		for _, ip := range ips {
			byteIps = append(byteIps, []byte(ip))
		}
		n, err := c.sessionService.CloseLocalSessionsByIp(ctx, byteIps, constant.SessionCloseStatus_DISCONNECTED_BY_ADMIN)
		if err != nil {
			return 0, err
		}
		count += n
	}

	return count, nil
}
