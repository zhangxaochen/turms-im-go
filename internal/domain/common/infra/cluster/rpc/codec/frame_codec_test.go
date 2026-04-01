package codec

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockReader wraps bytes.Reader to fulfill io.ByteReader
type mockReader struct {
	*bytes.Reader
}

func TestWriteAndReadFrame(t *testing.T) {
	// 1. Prepare data
	frame := RpcFrame{
		CodecID:   12,
		RequestID: 1001,
		Payload:   []byte("Hello, Turms RPC!"),
	}

	buf := new(bytes.Buffer)

	// 2. Encode
	err := WriteFrame(buf, frame)
	require.NoError(t, err)

	byteSlice := buf.Bytes()
	// Length should be Varint(len(Payload) + 6) + 6 + len(Payload)
	// payload size is 17. 17+6=23. Varint for 23 takes 1 byte. Total 24 bytes.
	assert.Equal(t, 24, len(byteSlice))

	// 3. Decode
	reader := mockReader{bytes.NewReader(byteSlice)}
	decodedFrame, err := ReadFrame(reader)
	require.NoError(t, err)

	assert.Equal(t, frame.CodecID, decodedFrame.CodecID)
	assert.Equal(t, frame.RequestID, decodedFrame.RequestID)
	assert.Equal(t, frame.Payload, decodedFrame.Payload)
}

func TestStickyPackets(t *testing.T) {
	frames := []RpcFrame{
		{CodecID: 1, RequestID: 100, Payload: []byte("P1")},
		{CodecID: 2, RequestID: 101, Payload: []byte("P2")},
		{CodecID: 3, RequestID: 102, Payload: []byte("P3")},
	}

	buf := new(bytes.Buffer)

	// Write all frames into one consecutive buffer (simulating sticky TCP packets)
	for _, frame := range frames {
		require.NoError(t, WriteFrame(buf, frame))
	}

	reader := mockReader{bytes.NewReader(buf.Bytes())}

	// Try reading them consecutively back
	for _, frame := range frames {
		decoded, err := ReadFrame(reader)
		require.NoError(t, err)

		assert.Equal(t, frame.CodecID, decoded.CodecID)
		assert.Equal(t, frame.RequestID, decoded.RequestID)
		assert.Equal(t, frame.Payload, decoded.Payload)
	}

	// Finally it should EOF
	_, err := ReadFrame(reader)
	assert.ErrorIs(t, err, io.EOF)
}

// stickyHalfReader simulates reading chunks simulating network delay
type stickyHalfReader struct {
	src      []byte
	pos      int
	chunks   []int // Size of pieces to serve simulating network
	chunkIdx int
}

func (r *stickyHalfReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.src) {
		return 0, io.EOF
	}

	// How many bytes to serve in this call limit?
	if r.chunkIdx >= len(r.chunks) {
		r.chunkIdx = 0
	}
	limit := r.chunks[r.chunkIdx]
	r.chunkIdx++

	remain := len(r.src) - r.pos
	if limit > remain {
		limit = remain
	}

	maxRead := len(p)
	if limit > maxRead {
		limit = maxRead
	}

	n = copy(p, r.src[r.pos:r.pos+limit])
	r.pos += n
	return n, nil
}

func (r *stickyHalfReader) ReadByte() (byte, error) {
	if r.pos >= len(r.src) {
		return 0, io.EOF
	}
	b := r.src[r.pos]
	r.pos++
	return b, nil
}

func TestHalfPackets(t *testing.T) {
	frame := RpcFrame{
		CodecID:   5,
		RequestID: 999,
		Payload:   bytes.Repeat([]byte("X"), 1000), // ~1KB packet
	}

	buf := new(bytes.Buffer)
	require.NoError(t, WriteFrame(buf, frame))

	// Simulate reading byte by byte or tiny chunks (extremely slow network)
	reader := &stickyHalfReader{
		src:    buf.Bytes(),
		chunks: []int{1, 2, 50, 100}, // Give small bites
	}

	decoded, err := ReadFrame(reader)
	require.NoError(t, err)

	assert.Equal(t, frame.CodecID, decoded.CodecID)
	assert.Equal(t, frame.RequestID, decoded.RequestID)
	assert.Equal(t, frame.Payload, decoded.Payload)
}

func TestInvalidRequestID(t *testing.T) {
	frame := RpcFrame{
		CodecID:   1,
		RequestID: -1,
		Payload:   []byte("P"),
	}

	err := WriteFrame(new(bytes.Buffer), frame)
	assert.ErrorIs(t, err, ErrInvalidRequestID)
}
