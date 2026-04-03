package udp

import (
	"context"
	"net"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

// @MappedFrom UdpNotificationType
type UdpNotificationType byte

const (
	HeartbeatNotification UdpNotificationType = iota
	GoOfflineNotification
)

// @MappedFrom UdpRequestType
type UdpRequestType byte

const (
	HeartbeatRequest UdpRequestType = iota
	GoOfflineRequest
)

// @MappedFrom UdpNotification
type UdpNotification struct {
	RecipientAddress net.Addr
	Type             UdpNotificationType
}

// @MappedFrom UdpNotification(InetSocketAddress recipientAddress, UdpNotificationType type)
func NewUdpNotification(recipientAddress net.Addr, notificationType UdpNotificationType) *UdpNotification {
	return &UdpNotification{
		RecipientAddress: recipientAddress,
		Type:             notificationType,
	}
}

// @MappedFrom UdpSignalRequest
type UdpSignalRequest struct {
	Type       UdpRequestType
	UserID     int64
	DeviceType protocol.DeviceType
	SessionID  int
}

// @MappedFrom UdpSignalRequest(UdpRequestType type, long userId, DeviceType deviceType, int sessionId)
func NewUdpSignalRequest(reqType UdpRequestType, userID int64, deviceType protocol.DeviceType, sessionID int) *UdpSignalRequest {
	return &UdpSignalRequest{
		Type:       reqType,
		UserID:     userID,
		DeviceType: deviceType,
		SessionID:  sessionID,
	}
}

// @MappedFrom parse(int number)
func ParseUdpRequestType(number int) UdpRequestType {
	return UdpRequestType(number)
}

// @MappedFrom getNumber()
func (t UdpRequestType) GetNumber() int {
	return int(t)
}

// @MappedFrom UdpRequestDispatcher
type UdpRequestDispatcher struct {
	sessionService   *session.SessionService
	notificationSink chan UdpNotification
	connection       *net.UDPConn
}

func NewUdpRequestDispatcher(sessionService *session.SessionService, enabled bool, host string, port int) *UdpRequestDispatcher {
	if !enabled {
		return &UdpRequestDispatcher{}
	}
	// Pending implementation: net.ListenUDP, run accept loop, manage sink
	return &UdpRequestDispatcher{
		sessionService:   sessionService,
		notificationSink: make(chan UdpNotification, 1024),
	}
}

// @MappedFrom sendSignal(InetSocketAddress address, UdpNotificationType signal)
func (d *UdpRequestDispatcher) SendSignal(address net.Addr, signal UdpNotificationType) {
	if d.notificationSink != nil {
		select {
		case d.notificationSink <- UdpNotification{
			RecipientAddress: address,
			Type:             signal,
		}:
		default:
			// Handle sink full
		}
	}
}

// @MappedFrom handleDatagramPackage(DatagramPacket packet)
func (d *UdpRequestDispatcher) HandleDatagramPackage(ctx context.Context, packet []byte, senderAddress net.Addr) error {
	req := d.ParseRequest(packet)
	if req == nil {
		return nil // MAP TO INVALID_REQUEST status code logic
	}

	switch req.Type {
	case HeartbeatRequest:
		// Pending implementation: sessionService.AuthAndUpdateHeartbeatTimestamp
		// And update UDP Address on session connection
	case GoOfflineRequest:
		// Pending implementation: sessionService.AuthAndCloseLocalSession
	}

	return nil
}

// @MappedFrom parseRequest(ByteBuf byteBuf)
func (d *UdpRequestDispatcher) ParseRequest(buffer []byte) *UdpSignalRequest {
	// Pending implementation: read UDP packet bytes
	return nil
}

// @MappedFrom get(ResponseStatusCode code)
func (d *UdpRequestDispatcher) GetBufferFromStatusCode(code constant.ResponseStatusCode) []byte {
	if code == constant.ResponseStatusCode_OK {
		return []byte{}
	}
	// Simplified mock: return 2 bytes representing the business code
	// Usually this would use encoding/binary byte order
	val := uint16(code)
	return []byte{byte(val >> 8), byte(val)}
}

// @MappedFrom get(UdpNotificationType type)
func (d *UdpRequestDispatcher) GetBufferFromNotificationType(t UdpNotificationType) []byte {
	// ordinal + 1 per the Java implementation
	return []byte{byte(t) + 1}
}
