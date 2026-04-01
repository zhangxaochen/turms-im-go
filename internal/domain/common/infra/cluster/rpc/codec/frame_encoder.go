package codec

import (
	"encoding/binary"
	"io"
)

// WriteFrame encodes and writes an RpcFrame to the provided io.Writer.
// The structure is:
// 1. Varint32 representing the total length of the frame payload
// 2. The payload itself:
//   - Codec ID (uint16, 2 bytes, BigEndian)
//   - Request ID (int32, 4 bytes, BigEndian)
//   - The actual byte data
func WriteFrame(w io.Writer, frame RpcFrame) error {
	if frame.RequestID < 0 {
		return ErrInvalidRequestID
	}

	payloadLen := HeaderSize + len(frame.Payload)

	// Encode length using Uvarint (up to 5 bytes for 32-bit uint)
	var lengthBuf [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(lengthBuf[:], uint64(payloadLen))

	// Write length prefix
	if _, err := w.Write(lengthBuf[:n]); err != nil {
		return err
	}

	// Prepare codec header
	var headerBuf [HeaderSize]byte
	binary.BigEndian.PutUint16(headerBuf[0:2], frame.CodecID)
	binary.BigEndian.PutUint32(headerBuf[2:6], uint32(frame.RequestID))

	// Write header
	if _, err := w.Write(headerBuf[:]); err != nil {
		return err
	}

	// Write payload
	if len(frame.Payload) > 0 {
		if _, err := w.Write(frame.Payload); err != nil {
			return err
		}
	}

	return nil
}
