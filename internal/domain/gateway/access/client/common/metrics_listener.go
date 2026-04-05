package common

import (
	"net"
	"sync/atomic"
)

// MetricsListener wraps a net.Listener to track connection-level metrics.
// @MappedFrom Java's .metrics(true, () -> new TurmsMicrometerChannelMetricsRecorder(...))
// In Java, Micrometer records: connections active, bytes received/sent, connect time, etc.
// This Go equivalent tracks the same core metrics using atomics for lock-free access.
type MetricsListener struct {
	net.Listener
	metricsName       string
	ActiveConnections atomic.Int64
	TotalConnections  atomic.Int64
	TotalBytesRead    atomic.Int64
	TotalBytesWritten atomic.Int64
}

func NewMetricsListener(l net.Listener, metricsName string) *MetricsListener {
	return &MetricsListener{
		Listener:    l,
		metricsName: metricsName,
	}
}

func (ml *MetricsListener) Accept() (net.Conn, error) {
	conn, err := ml.Listener.Accept()
	if err != nil {
		return nil, err
	}
	ml.ActiveConnections.Add(1)
	ml.TotalConnections.Add(1)
	return &MetricsConn{
		Conn:     conn,
		listener: ml,
	}, nil
}

// MetricsName returns the name used for this metrics recorder.
// @MappedFrom MetricNameConst.TURMS_GATEWAY_SERVER_TCP / TURMS_GATEWAY_SERVER_WS
func (ml *MetricsListener) MetricsName() string {
	return ml.metricsName
}

// Snapshot returns current metric values for external consumption.
func (ml *MetricsListener) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		ActiveConnections: ml.ActiveConnections.Load(),
		TotalConnections:  ml.TotalConnections.Load(),
		TotalBytesRead:    ml.TotalBytesRead.Load(),
		TotalBytesWritten: ml.TotalBytesWritten.Load(),
	}
}

type MetricsSnapshot struct {
	ActiveConnections int64
	TotalConnections  int64
	TotalBytesRead    int64
	TotalBytesWritten int64
}

// MetricsConn wraps a net.Conn to track bytes read/written and connection lifecycle.
type MetricsConn struct {
	net.Conn
	listener *MetricsListener
	closed   atomic.Bool
}

func (mc *MetricsConn) Read(b []byte) (n int, err error) {
	n, err = mc.Conn.Read(b)
	if n > 0 {
		mc.listener.TotalBytesRead.Add(int64(n))
	}
	return
}

func (mc *MetricsConn) Write(b []byte) (n int, err error) {
	n, err = mc.Conn.Write(b)
	if n > 0 {
		mc.listener.TotalBytesWritten.Add(int64(n))
	}
	return
}

func (mc *MetricsConn) Close() error {
	if mc.closed.CompareAndSwap(false, true) {
		mc.listener.ActiveConnections.Add(-1)
	}
	return mc.Conn.Close()
}
