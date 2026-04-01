package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/domain/common/infra/cluster/rpc/codec"
)

func TestRouter_Dispatch(t *testing.T) {
	router := NewRouter()

	// Register a mock handler
	const codecID uint16 = 404
	router.Register(codecID, func(ctx context.Context, payload []byte) ([]byte, error) {
		assert.Equal(t, []byte("request-data"), payload)
		return []byte("response-data"), nil
	})

	// Test Successful Dispatch
	frame := &codec.RpcFrame{
		CodecID: codecID,
		Payload: []byte("request-data"),
	}

	res, err := router.Dispatch(context.Background(), frame)
	require.NoError(t, err)
	assert.Equal(t, []byte("response-data"), res)

	// Test Unregistered Dispatch
	unregisteredFrame := &codec.RpcFrame{
		CodecID: 999, // Unregistered
		Payload: []byte("who-are-you"),
	}
	_, err = router.Dispatch(context.Background(), unregisteredFrame)
	assert.ErrorIs(t, err, ErrHandlerNotFound)
}
