package codec

import "errors"

const (
	// HeaderSize is the fixed size of CodecID (2) + RequestID (4)
	HeaderSize = 6
)

var (
	ErrBufferTooSmall   = errors.New("buffer is too small to parse the header")
	ErrInvalidRequestID = errors.New("request ID must be greater than 0")
	ErrPacketTooLarge   = errors.New("packet length exceeds maximum allowed")
)

// RpcFrame represents the decoded or to-be-encoded frame data structure
type RpcFrame struct {
	CodecID   uint16
	RequestID int32
	Payload   []byte // The inner payload
}
