package common

import (
	"log"
	"net"
	"sync/atomic"
)

// WiretapListener wraps a net.Listener to log all accepted connections.
// @MappedFrom Reactor Netty's .wiretap() configuration
// In Java, wiretap enables DEBUG-level logging of all inbound/outbound channel events.
// In Go, this wraps connections with a logging proxy that logs read/write data.
type WiretapListener struct {
	net.Listener
}

func NewWiretapListener(l net.Listener) *WiretapListener {
	return &WiretapListener{Listener: l}
}

func (wl *WiretapListener) Accept() (net.Conn, error) {
	conn, err := wl.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return NewWiretapConn(conn), nil
}

// WiretapConn wraps a net.Conn to log read/write activity at DEBUG level.
type WiretapConn struct {
	net.Conn
	id          uint64
	bytesRead   atomic.Int64
	bytesWritten atomic.Int64
}

var wiretapConnCounter atomic.Uint64

func NewWiretapConn(conn net.Conn) *WiretapConn {
	id := wiretapConnCounter.Add(1)
	log.Printf("[wiretap] connection #%d accepted from %s", id, conn.RemoteAddr())
	return &WiretapConn{
		Conn: conn,
		id:   id,
	}
}

func (wc *WiretapConn) Read(b []byte) (n int, err error) {
	n, err = wc.Conn.Read(b)
	if n > 0 {
		wc.bytesRead.Add(int64(n))
		log.Printf("[wiretap] #%d READ %d bytes (total: %d) from %s",
			wc.id, n, wc.bytesRead.Load(), wc.RemoteAddr())
	}
	return
}

func (wc *WiretapConn) Write(b []byte) (n int, err error) {
	n, err = wc.Conn.Write(b)
	if n > 0 {
		wc.bytesWritten.Add(int64(n))
		log.Printf("[wiretap] #%d WRITE %d bytes (total: %d) to %s",
			wc.id, n, wc.bytesWritten.Load(), wc.RemoteAddr())
	}
	return
}

func (wc *WiretapConn) Close() error {
	log.Printf("[wiretap] #%d CLOSED (read: %d, written: %d) %s",
		wc.id, wc.bytesRead.Load(), wc.bytesWritten.Load(), wc.RemoteAddr())
	return wc.Conn.Close()
}
