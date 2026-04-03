package rpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"im.turms/server/internal/domain/common/infra/cluster/connection"
	"im.turms/server/internal/domain/common/infra/cluster/discovery"
)

var (
	ErrMemberNotFound     = errors.New("healthy member not found")
	ErrConnectionNotFound = errors.New("connection not found")
	ErrRequestTimeout     = errors.New("rpc request timeout")
	ErrNotAllowedToSend   = errors.New("current node type is not allowed to send this request")
)

type RpcService struct {
	nodeType          discovery.NodeType
	discoveryService  *discovery.DiscoveryService
	connectionService *connection.ConnectionService
	router            *Router

	// defaultRequestTimeoutDuration represents the fallback timeout
	defaultRequestTimeoutDuration time.Duration

	// pendingRequests maps RequestID (int32) to chan RpcResponse
	pendingRequests sync.Map
	requestIDSeq    int32
}

func NewRpcService(
	nodeType discovery.NodeType,
	discoveryService *discovery.DiscoveryService,
	connectionService *connection.ConnectionService,
	router *Router,
	timeoutMillis int,
) *RpcService {
	return &RpcService{
		nodeType:                      nodeType,
		discoveryService:              discoveryService,
		connectionService:             connectionService,
		router:                        router,
		defaultRequestTimeoutDuration: time.Duration(timeoutMillis) * time.Millisecond,
		requestIDSeq:                  0,
	}
}

// RequestResponse routes the request to an appropriate member node and waits for a response.
// If memberNodeId is provided, it specifically targets that node.
func (s *RpcService) RequestResponse(ctx context.Context, memberNodeId string, request RpcRequest) (*RpcResponse, error) {
	if err := s.assertCurrentNodeIsAllowedToSend(request); err != nil {
		return nil, err
	}

	if memberNodeId == "" {
		member, err := s.selectTargetMember(request)
		if err != nil {
			return nil, err
		}
		memberNodeId = member.NodeID
	}

	// Local dispatch check
	if memberNodeId == s.discoveryService.GetLocalNodeID() {
		// TODO: Local dispatch via Router.
		// return s.router.DispatchLocal(ctx, request)
		return nil, errors.New("local dispatch not yet fully implemented")
	}

	return s.requestResponseNode(ctx, memberNodeId, request)
}

// RequestResponsesFromOtherMembers sends the request to all active connected members of the specified node type.
func (s *RpcService) RequestResponsesFromOtherMembers(ctx context.Context, request RpcRequest, rejectIfMissingAnyConnection bool) (map[string]*RpcResponse, error) {
	if err := s.assertCurrentNodeIsAllowedToSend(request); err != nil {
		return nil, err
	}

	members := s.getOtherActiveConnectedMembersToRespond(request)
	if len(members) == 0 {
		return nil, ErrMemberNotFound
	}

	if rejectIfMissingAnyConnection && !s.connectionService.IsHasConnectedToAllMembers() {
		return nil, fmt.Errorf("%w: not all connections are established", ErrConnectionNotFound)
	}

	responses := make(map[string]*RpcResponse)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, member := range members {
		wg.Add(1)
		go func(nodeID string) {
			defer wg.Done()
			resp, err := s.requestResponseNode(ctx, nodeID, request)
			if err != nil {
				// Record error response or skip. For now, wrap error in RpcResponse.
				resp = &RpcResponse{Err: err}
			}
			mu.Lock()
			responses[nodeID] = resp
			mu.Unlock()
		}(member.NodeID)
	}

	wg.Wait()
	return responses, nil
}

// requestResponseNode handles the actual encoding, sending, and waiting.
func (s *RpcService) requestResponseNode(ctx context.Context, memberNodeId string, request RpcRequest) (*RpcResponse, error) {
	// 1. Check connection
	/*
		conn := s.connectionService.GetMemberConnection(memberNodeId)
		if conn == nil {
			return nil, ErrConnectionNotFound
		}
	*/

	// 2. Generate Request ID
	// _ = atomic.AddInt32(&s.requestIDSeq, 1)

	// 3. Register future
	/*
		respChan := make(chan *RpcResponse, 1)
		s.pendingRequests.Store(reqID, respChan)
		defer s.pendingRequests.Delete(reqID)
	*/

	// 4. Encode Payload
	/*
		payload, err := request.Payload()
		frame := codec.RpcFrame{CodecID: request.CodecID(), RequestID: reqID, Payload: payload}
		encodedBody := s.codecService.Serialize(frame)
	*/

	// 5. Build and send frame
	// conn.Send(encodedBody)

	// TODO: 完整的 TCP/UDP 底层 P2P 网络请求发送逻辑暂未实现。
	// 请参考 infrastructure 规划 - 后续可能会使用 gRPC 替代，或引入标准 go net 建立 P2P peer 池。
	return nil, errors.New("network send is stubbed // TODO: implement binary RpcFrame TCP write or gRPC proxy")
}

// Helper functions for validating node types and selecting members.

func (s *RpcService) assertCurrentNodeIsAllowedToSend(request RpcRequest) error {
	reqType := request.NodeTypeToRequest()
	if reqType == NodeTypeToHandleRpcBoth {
		return nil
	}
	if reqType.MapToDiscoveryNodeType() != s.nodeType {
		return fmt.Errorf("%w: current %s cannot send %s which requires %s", ErrNotAllowedToSend, s.nodeType, request.Name(), reqType)
	}
	return nil
}

func (s *RpcService) getOtherActiveConnectedMembersToRespond(request RpcRequest) []*discovery.Member {
	targetNodeType := request.NodeTypeToRespond().MapToDiscoveryNodeType()
	allMembers := s.discoveryService.GetMembers()

	var validMembers []*discovery.Member
	for _, m := range allMembers {
		// // TODO: Filter out local node
		// if m.NodeID == s.discoveryService.LocalNodeID() { continue }

		if m.IsActive && m.IsHealthy {
			if targetNodeType == "" || m.NodeType == targetNodeType {
				validMembers = append(validMembers, m)
			}
		}
	}
	return validMembers
}

func (s *RpcService) selectTargetMember(request RpcRequest) (*discovery.Member, error) {
	members := s.getOtherActiveConnectedMembersToRespond(request)
	if len(members) == 0 {
		return nil, ErrMemberNotFound
	}

	// TODO: Add proper load balancing (e.g. random, round robin).
	// Currently returning the first one.
	return members[0], nil
}
