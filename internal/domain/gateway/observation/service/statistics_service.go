package service

import (
	"im.turms/server/internal/domain/gateway/session"
)

type StatisticsService struct {
	sessionService *session.SessionService
}

func NewStatisticsService(sessionService *session.SessionService) *StatisticsService {
	return &StatisticsService{
		sessionService: sessionService,
	}
}

// CountLocalOnlineUsers returns the approximate number of active connections in this gateway node.
// @MappedFrom countLocalOnlineUsers()
func (s *StatisticsService) CountLocalOnlineUsers() int {
	return s.sessionService.CountOnlineUsers()
}
