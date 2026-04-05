package common

import (
	"net"

	"im.turms/server/internal/domain/gateway/session"
)

// UserSessionWrapper binds the network connection and the user session together from the perspective of the access layer.
type UserSessionWrapper struct {
	IP          string
	IPStr       string
	UserSession          *session.UserSession
	Conn                 session.Connection
	OnSessionEstablished func(*session.UserSession)
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
