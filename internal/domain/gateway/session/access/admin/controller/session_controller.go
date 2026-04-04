package controller

import (
	"context"
	"net"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session"
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
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

// QuerySessions returns session info based on the provided user IDs.
// @MappedFrom querySessions(@QueryParam(required = false) Set<Long> ids)
func (c *SessionController) QuerySessions(ctx context.Context, ids []int64) ([]*sessionbo.UserSessionsInfo, error) {
	return c.sessionService.GetLocalUserSessionsInfo(ctx, ids), nil
}

// DeleteSessions deletes sessions based on the provided user IDs and/or IPs.
// @MappedFrom deleteSessions(@QueryParam(required = false) Set<Long> ids, @QueryParam(required = false) Set<String> ips)
func (c *SessionController) DeleteSessions(ctx context.Context, ids []int64, ips []string) (int, error) {
	closeReason := sessionbo.NewCloseReason(constant.SessionCloseStatus_DISCONNECTED_BY_ADMIN)
	if len(ids) == 0 && len(ips) == 0 {
		return c.sessionService.CloseAllLocalSessions(ctx, closeReason)
	}

	var ipsBytes [][]byte
	if len(ips) > 0 {
		ipsBytes = make([][]byte, len(ips))
		for i, ip := range ips {
			parsedIp := net.ParseIP(ip)
			if parsedIp == nil {
				ipsBytes[i] = []byte(ip)
			} else {
				if ipv4 := parsedIp.To4(); ipv4 != nil {
					ipsBytes[i] = ipv4
				} else {
					ipsBytes[i] = parsedIp
				}
			}
		}
	}
	return c.sessionService.CloseLocalSessions(ctx, ids, ipsBytes, closeReason)
}
