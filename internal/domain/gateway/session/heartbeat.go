package session

import (
	"context"
	"log"
	"time"

	"im.turms/server/pkg/protocol"
)

// HeartbeatManager is responsible for kicking offline users when their heartbeats timeout,
// and potentially synchronizing local online statues with Redis.
type HeartbeatManager struct {
	sessionService *SessionService
	interval       time.Duration
	timeout        time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

func NewHeartbeatManager(ss *SessionService, interval, timeout time.Duration) *HeartbeatManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &HeartbeatManager{
		sessionService: ss,
		interval:       interval,
		timeout:        timeout,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (m *HeartbeatManager) Start() {
	go m.cleanupLoop()
}

func (m *HeartbeatManager) Stop() {
	m.cancel()
}

func (m *HeartbeatManager) cleanupLoop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.checkHeartbeats()
			m.syncWithGlobal()
		}
	}
}

func (m *HeartbeatManager) checkHeartbeats() {
	now := time.Now().UnixMilli()
	timeoutMs := m.timeout.Milliseconds()

	smap := m.sessionService.shardedMap
	for _, shard := range smap.shards {
		// Read lock first to find dead connections
		shard.RLock()

		// To avoid holding locks for too long, we collect dead sessions
		type deadSession struct {
			userID     int64
			deviceType protocol.DeviceType
			conn       Connection
		}
		var toRemove []deadSession

		for userID, manager := range shard.m {
			manager.mu.RLock()
			for deviceType, userSession := range manager.Sessions {
				lastActivity := userSession.GetLastHeartbeatRequestTimestamp()
				if now-lastActivity > timeoutMs {
					toRemove = append(toRemove, deadSession{
						userID:     userID,
						deviceType: deviceType,
						conn:       userSession.Conn,
					})
				}
			}
			manager.mu.RUnlock()
		}
		shard.RUnlock()

		// Now we kick them using UnregisterSession
		for _, s := range toRemove {
			log.Printf("Session %d:%s kicked due to heartbeat timeout", s.userID, s.deviceType.String())
			m.sessionService.UnregisterSession(s.userID, s.deviceType, s.conn)
		}
	}
}

func (m *HeartbeatManager) syncWithGlobal() {
	// TODO: Phase 4/5, bulk push active local user IDs to Redis to let
	// other nodes know this node holds the active connection.
}
