package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	turmsredis "im.turms/server/internal/storage/redis"
)

type NodeType string

const (
	NodeTypeService NodeType = "SERVICE"
	NodeTypeGateway NodeType = "GATEWAY"
	NodeTypeAdmin   NodeType = "ADMIN"
)

type Member struct {
	ClusterID        string    `json:"cluster_id"`
	NodeID           string    `json:"node_id"`
	Zone             string    `json:"zone"`
	Name             string    `json:"name"`
	NodeType         NodeType  `json:"node_type"`
	IsSeed           bool      `json:"is_seed"`
	IsLeaderEligible bool      `json:"is_leader_eligible"`
	Priority         int       `json:"priority"`
	MemberHost       string    `json:"member_host"`
	MemberPort       int       `json:"member_port"`
	AdminAPIAddress  string    `json:"admin_api_address"`
	WsAddress        string    `json:"ws_address"`
	TcpAddress       string    `json:"tcp_address"`
	UdpAddress       string    `json:"udp_address"`
	IsActive         bool      `json:"is_active"`
	IsHealthy        bool      `json:"is_healthy"`
	LastHeartbeat    time.Time `json:"last_heartbeat"`
}

type DiscoveryService struct {
	redisClient *turmsredis.Client
	localMember *Member

	allKnownMembers sync.Map // map[string]*Member string=node_id

	heartbeatInterval time.Duration
	ttl               time.Duration

	leaderID string
	leaderMu sync.RWMutex

	stopCtx    context.Context
	stopCancel context.CancelFunc
	wg         sync.WaitGroup

	listeners []MembersChangeListener
	lisMu     sync.RWMutex
}

type MembersChangeListener interface {
	OnMembersChange()
}

func NewDiscoveryService(redisClient *turmsredis.Client, localMember *Member, heartbeatIntervalSeconds int) *DiscoveryService {
	ctx, cancel := context.WithCancel(context.Background())
	interval := time.Duration(heartbeatIntervalSeconds) * time.Second
	return &DiscoveryService{
		redisClient:       redisClient,
		localMember:       localMember,
		heartbeatInterval: interval,
		ttl:               interval * 3, // Missing 3 heartbeats means offline
		stopCtx:           ctx,
		stopCancel:        cancel,
	}
}

func (s *DiscoveryService) Start() error {
	// 1. Initial sync
	err := s.syncMembers()
	if err != nil {
		return fmt.Errorf("initial sync members failed: %w", err)
	}

	// 2. Register local node
	err = s.registerLocalMember()
	if err != nil {
		return fmt.Errorf("registering local member failed: %w", err)
	}

	// 3. Initial Leader Election
	if s.localMember.IsLeaderEligible {
		s.runLeaderElection()
	}

	// 4. Start background routines
	s.wg.Add(3)
	go s.heartbeatRoutine()
	go s.leaderElectionRoutine()
	go s.syncRoutine()

	return nil
}

func (s *DiscoveryService) Stop() {
	s.stopCancel()
	s.wg.Wait()
	// Unregister gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.redisClient.RDB.Del(ctx, s.getMemberKeyString(s.localMember.NodeID))
}

func (s *DiscoveryService) AddListener(listener MembersChangeListener) {
	s.lisMu.Lock()
	defer s.lisMu.Unlock()
	s.listeners = append(s.listeners, listener)
}

func (s *DiscoveryService) registerLocalMember() error {
	s.localMember.LastHeartbeat = time.Now()
	data, err := json.Marshal(s.localMember)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = s.redisClient.RDB.Set(ctx, s.getMemberKeyString(s.localMember.NodeID), data, s.ttl).Err()
	if err == nil {
		s.allKnownMembers.Store(s.localMember.NodeID, s.localMember)
	}
	return err
}

func (s *DiscoveryService) heartbeatRoutine() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCtx.Done():
			return
		case <-ticker.C:
			s.localMember.LastHeartbeat = time.Now()
			data, _ := json.Marshal(s.localMember)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			s.redisClient.RDB.Set(ctx, s.getMemberKeyString(s.localMember.NodeID), data, s.ttl)
			cancel()
		}
	}
}

func (s *DiscoveryService) syncRoutine() {
	defer s.wg.Done()
	// Sync slightly faster than heartbeat to catch offline nodes quickly
	ticker := time.NewTicker(s.heartbeatInterval / 2)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCtx.Done():
			return
		case <-ticker.C:
			_ = s.syncMembers()
		}
	}
}

func (s *DiscoveryService) syncMembers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchPattern := fmt.Sprintf("turms:cluster:%s:members:*", s.localMember.ClusterID)
	// For production with thousand nodes, SCAN is required.
	// For small clusters, KEYS is functionally identical but SCAN is safer.
	var cursor uint64
	var allKeys []string
	for {
		keys, nextCursor, err := s.redisClient.RDB.Scan(ctx, cursor, matchPattern, 100).Result()
		if err != nil {
			return err
		}
		allKeys = append(allKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(allKeys) == 0 {
		s.allKnownMembers.Range(func(key, value any) bool {
			s.allKnownMembers.Delete(key)
			return true
		})
		s.notifyListeners()
		return nil
	}

	// Use MGET to fetch all members atomically
	values, err := s.redisClient.RDB.MGet(ctx, allKeys...).Result()
	if err != nil {
		return err
	}

	newMembers := make(map[string]*Member)
	for _, v := range values {
		if v == nil {
			continue
		}
		strVal, ok := v.(string)
		if !ok {
			continue
		}
		var m Member
		if err := json.Unmarshal([]byte(strVal), &m); err == nil {
			newMembers[m.NodeID] = &m
			s.allKnownMembers.Store(m.NodeID, &m)
		}
	}

	// Remove stale members that expired and vanished from Redis
	var staleNodes []string
	s.allKnownMembers.Range(func(key, value any) bool {
		nodeID := key.(string)
		if _, exists := newMembers[nodeID]; !exists {
			staleNodes = append(staleNodes, nodeID)
		}
		return true
	})

	for _, nodeID := range staleNodes {
		s.allKnownMembers.Delete(nodeID)
	}

	// We can add logic to compare diffs and notify IF something actually changed
	s.notifyListeners()
	return nil
}

func (s *DiscoveryService) runLeaderElection() {
	if !s.localMember.IsLeaderEligible {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	leaderKey := fmt.Sprintf("turms:cluster:%s:leader", s.localMember.ClusterID)

	ok, err := s.redisClient.RDB.SetNX(ctx, leaderKey, s.localMember.NodeID, s.ttl).Result()
	if err != nil {
		return
	}
	if ok {
		s.leaderMu.Lock()
		s.leaderID = s.localMember.NodeID
		s.leaderMu.Unlock()
	} else {
		currentLeader, _ := s.redisClient.RDB.Get(ctx, leaderKey).Result()
		if currentLeader == s.localMember.NodeID {
			s.redisClient.RDB.SetXX(ctx, leaderKey, s.localMember.NodeID, s.ttl)
		}
		s.leaderMu.Lock()
		s.leaderID = currentLeader
		s.leaderMu.Unlock()
	}
}

func (s *DiscoveryService) leaderElectionRoutine() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()

	// Initial run is already done in Start()

	for {
		select {
		case <-s.stopCtx.Done():
			return
		case <-ticker.C:
			s.runLeaderElection()
		}
	}
}

func (s *DiscoveryService) notifyListeners() {
	s.lisMu.RLock()
	defer s.lisMu.RUnlock()
	for _, l := range s.listeners {
		l.OnMembersChange()
	}
}

// GetMembers returns all active members
func (s *DiscoveryService) GetMembers() []*Member {
	var members []*Member
	s.allKnownMembers.Range(func(key, value any) bool {
		members = append(members, value.(*Member))
		return true
	})
	return members
}

func (s *DiscoveryService) GetLocalNodeID() string {
	if s.localMember != nil {
		return s.localMember.NodeID
	}
	return ""
}

func (s *DiscoveryService) GetLeaderID() string {
	s.leaderMu.RLock()
	defer s.leaderMu.RUnlock()
	return s.leaderID
}

func (s *DiscoveryService) getMemberKeyString(nodeID string) string {
	return fmt.Sprintf("turms:cluster:%s:members:%s", s.localMember.ClusterID, nodeID)
}
