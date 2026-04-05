package common

import (
	"net"
	"sync/atomic"
	"time"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session"
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
)

// UserSessionWrapper binds the network connection and the user session together from the perspective of the access layer.
type UserSessionWrapper struct {
	IP          string
	IPStr       string
	UserSession          *session.UserSession
	Conn                 session.Connection
	OnSessionEstablished func(*session.UserSession)

	establishTimeoutTimer *time.Timer
	establishTimeoutFired atomic.Bool
}

// @MappedFrom UserSessionWrapper(NetConnection connection, InetSocketAddress address, int establishTimeoutMillis, Consumer<UserSession> onSessionEstablished)
func NewUserSessionWrapper(conn session.Connection, addr net.Addr, establishTimeoutMillis int, onSessionEstablished func(*session.UserSession)) *UserSessionWrapper {
	w := &UserSessionWrapper{
		Conn:                 conn,
		OnSessionEstablished: onSessionEstablished,
	}
	if establishTimeoutMillis > 0 {
		w.establishTimeoutTimer = time.AfterFunc(time.Duration(establishTimeoutMillis)*time.Millisecond, func() {
			if w.UserSession == nil || !w.UserSession.IsOpen() {
				w.establishTimeoutFired.Store(true)
				w.Conn.CloseWithReason(sessionbo.NewCloseReason(constant.SessionCloseStatus_LOGIN_TIMEOUT))
			}
		})
	}
	return w
}

func (w *UserSessionWrapper) CancelEstablishTimeout() bool {
	if w.establishTimeoutTimer != nil {
		w.establishTimeoutTimer.Stop()
		return !w.establishTimeoutFired.Load()
	}
	return true
}

func (w *UserSessionWrapper) GetConnection() session.Connection {
	return w.Conn
}

// @MappedFrom getIp()
func (w *UserSessionWrapper) GetIP() string {
	if w.IP == "" && w.Conn != nil {
		if addr := w.Conn.GetAddress(); addr != nil {
			if tcpAddr, ok := addr.(*net.TCPAddr); ok {
				w.IP = string(tcpAddr.IP)
			} else if udpAddr, ok := addr.(*net.UDPAddr); ok {
				w.IP = string(udpAddr.IP)
			}
		}
	}
	return w.IP
}

// @MappedFrom getIpStr()
func (w *UserSessionWrapper) GetIPStr() string {
	if w.IPStr == "" && w.Conn != nil {
		if addr := w.Conn.GetAddress(); addr != nil {
			if tcpAddr, ok := addr.(*net.TCPAddr); ok {
				w.IPStr = tcpAddr.IP.String()
			} else if udpAddr, ok := addr.(*net.UDPAddr); ok {
				w.IPStr = udpAddr.IP.String()
			}
		}
	}
	return w.IPStr
}

// @MappedFrom setUserSession(UserSession userSession)
func (w *UserSessionWrapper) SetUserSession(userSession *session.UserSession) {
	w.UserSession = userSession
	w.UserSession.SetConnection(w.Conn, net.IP(w.GetIP()))
	if w.OnSessionEstablished != nil {
		w.OnSessionEstablished(userSession)
	}
}

// @MappedFrom hasUserSession()
func (w *UserSessionWrapper) HasUserSession() bool {
	return w.UserSession != nil
}
