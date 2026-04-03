package common

import (
	"im.turms/server/internal/domain/gateway/session"
)

// UserSessionWrapper binds the network connection and the user session together from the perspective of the access layer.
type UserSessionWrapper struct {
	IP          string
	IPStr       string
	UserSession *session.UserSession
}

// @MappedFrom getIp()
func (w *UserSessionWrapper) GetIP() string {
	return w.IP
}

// @MappedFrom getIpStr()
func (w *UserSessionWrapper) GetIPStr() string {
	return w.IPStr
}

// @MappedFrom setUserSession(UserSession userSession)
func (w *UserSessionWrapper) SetUserSession(userSession *session.UserSession) {
	w.UserSession = userSession
}

// @MappedFrom hasUserSession()
func (w *UserSessionWrapper) HasUserSession() bool {
	return w.UserSession != nil
}
