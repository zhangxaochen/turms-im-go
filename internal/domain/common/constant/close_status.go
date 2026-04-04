package constant

// SessionCloseStatus defines reasons for session closure
type SessionCloseStatus int32

const (
	SessionCloseStatus_ILLEGAL_REQUEST                    SessionCloseStatus = 100
	SessionCloseStatus_SERVER_ERROR                       SessionCloseStatus = 101
	SessionCloseStatus_SERVER_CLOSED                      SessionCloseStatus = 102
	SessionCloseStatus_SERVER_UNAVAILABLE                 SessionCloseStatus = 103
	SessionCloseStatus_CONNECTION_CLOSED                  SessionCloseStatus = 104
	SessionCloseStatus_UNKNOWN_ERROR                      SessionCloseStatus = 105
	SessionCloseStatus_DISCONNECTED_BY_CLIENT             SessionCloseStatus = 106
	SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE       SessionCloseStatus = 107
	SessionCloseStatus_DISCONNECTED_BY_ADMIN              SessionCloseStatus = 108
	SessionCloseStatus_USER_IS_DELETED_OR_INACTIVE        SessionCloseStatus = 109
	SessionCloseStatus_HEARTBEAT_TIMEOUT                  SessionCloseStatus = 110
	SessionCloseStatus_LOGIN_TIMEOUT                      SessionCloseStatus = 111
	SessionCloseStatus_SWITCH                             SessionCloseStatus = 112
	SessionCloseStatus_DISCONNECTED_BY_CLIENT_REDUNDANTLY SessionCloseStatus = 113
	SessionCloseStatus_DISCONNECTED_BY_SERVER             SessionCloseStatus = 114
	SessionCloseStatus_REDUNDANT_REQUEST                  SessionCloseStatus = 115
)
