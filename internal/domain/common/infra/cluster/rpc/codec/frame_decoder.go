package codec

import (
	"encoding/binary"
	"io"
)

// ReadFrame decodes an RpcFrame from an io.ByteReader (which must also be an io.Reader).
// It blocks until a full frame is read or an error occurs (such as io.EOF).
// This naturally handles strict packing (half packets/sticky packets) as long as the
// underlying reader buffers them properly.
func ReadFrame(r interface {
	io.Reader
	io.ByteReader
}) (*RpcFrame, error) {
	// 1. Read the length Varint prefix from the stream
	payloadLen, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	// Optional check against a max packet size (e.g., 8MB) to prevent OOM
	// if payloadLen > MaxPacketSize {
	//    return nil, ErrPacketTooLarge
	// }

	if payloadLen < HeaderSize {
		return nil, ErrBufferTooSmall
	}

	// 2. Read exactly `payloadLen` bytes
	fullPayload := make([]byte, payloadLen)
	if _, err := io.ReadFull(r, fullPayload); err != nil {
		return nil, err
	}

	// 3. Extract the components
	codecID := binary.BigEndian.Uint16(fullPayload[0:2])
	requestID := int32(binary.BigEndian.Uint32(fullPayload[2:6]))

	if requestID < 0 {
		return nil, ErrInvalidRequestID
	}

	frame := &RpcFrame{
		CodecID:   codecID,
		RequestID: requestID,
	}

	// The remaining bytes form the body payload
	if payloadLen > HeaderSize {
		frame.Payload = fullPayload[HeaderSize:]
	}

	return frame, nil
}
