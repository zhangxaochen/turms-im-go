package exception

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
)

// IsDisconnectedClientError checks if the error represents a client disconnection,
// similar to Java's ThrowableUtil.isDisconnectedClientError.
func IsDisconnectedClientError(err error) bool {
	if err == nil {
		return false
	}
	// Direct matches for common EOF/Closed errors
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	// Check net.OpError (e.g., read/write on a closed socket)
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if errors.Is(opErr.Err, syscall.ECONNRESET) ||
			errors.Is(opErr.Err, syscall.EPIPE) ||
			errors.Is(opErr.Err, syscall.ESHUTDOWN) {
			return true
		}
	}
	// Check os.SyscallError
	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		if errors.Is(syscallErr.Err, syscall.ECONNRESET) ||
			errors.Is(syscallErr.Err, syscall.EPIPE) ||
			errors.Is(syscallErr.Err, syscall.ESHUTDOWN) {
			return true
		}
	}
	// Fallback to string matching for platform-dependent messages
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "use of closed network connection") ||
		strings.Contains(msg, "connection was forcibly closed")
}
