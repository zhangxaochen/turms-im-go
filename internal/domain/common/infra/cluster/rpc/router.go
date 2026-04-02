package rpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"im.turms/server/internal/domain/common/infra/cluster/rpc/codec"
)

var (
	ErrHandlerNotFound = errors.New("rpc handler not found for codec id")
)

// HandlerFunc defines the signature for processing an incoming RPC request payload.
// It receives the byte payload of the RpcFrame, and should return the correctly encoded response payload bytes if any.
type HandlerFunc func(ctx context.Context, payload []byte) ([]byte, error)

// Router handles the registration and dispatching of RPC messages based on CodecID.
type Router struct {
	handlers map[uint16]HandlerFunc
	mu       sync.RWMutex
}

// NewRouter creates a new RPC message dispatch router.
func NewRouter() *Router {
	return &Router{
		handlers: make(map[uint16]HandlerFunc),
	}
}

// Register registers a handler function for a specific CodecID.
func (r *Router) Register(codecID uint16, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[codecID] = handler
}

// Dispatch routes the incoming RPC frame to its corresponding handler.
// It returns an error if the codec ID is not registered, or if the handler fails.
// @MappedFrom dispatch(TracingContext context, ServiceRequest serviceRequest)
func (r *Router) Dispatch(ctx context.Context, frame *codec.RpcFrame) ([]byte, error) {
	r.mu.RLock()
	handler, exists := r.handlers[frame.CodecID]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %d", ErrHandlerNotFound, frame.CodecID)
	}

	return handler(ctx, frame.Payload)
}
