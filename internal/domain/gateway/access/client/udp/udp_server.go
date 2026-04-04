package udp

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

// @MappedFrom UdpNotificationType
type UdpNotificationType byte

const (
	OpenConnectionNotification UdpNotificationType = iota // mapped to Java's OPEN_CONNECTION
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
	SessionID  int64
}

// @MappedFrom UdpSignalRequest(UdpRequestType type, long userId, DeviceType deviceType, int sessionId)
func NewUdpSignalRequest(reqType UdpRequestType, userID int64, deviceType protocol.DeviceType, sessionID int64) *UdpSignalRequest {
	return &UdpSignalRequest{
		Type:       reqType,
		UserID:     userID,
		DeviceType: deviceType,
		SessionID:  sessionID,
	}
}

// @MappedFrom parse(int number)
func ParseUdpRequestType(number int) (UdpRequestType, error) {
	index := number - 1
	if index >= 0 && index <= 1 {
		return UdpRequestType(index), nil
	}
	return 0, fmt.Errorf("invalid UDP request type number: %d", number)
}

// @MappedFrom getNumber()
func (t UdpRequestType) GetNumber() int {
	return int(t) + 1
}

var (
	Instance *UdpRequestDispatcher
)

// @MappedFrom UdpRequestDispatcher
type UdpRequestDispatcher struct {
	sessionService   *session.SessionService
	notificationSink chan UdpNotification
	connection       *net.UDPConn
	statusPool       sync.Map // Added for caching ResponseStatusCode buffers
	stopChan         chan struct{}
}

func NewUdpRequestDispatcher(sessionService *session.SessionService, enabled bool, host string, port int) *UdpRequestDispatcher {
	d := &UdpRequestDispatcher{
		sessionService: sessionService,
		stopChan:       make(chan struct{}),
	}
	Instance = d

	if !enabled {
		return d
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v", err)
		return d
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Failed to listen UDP: %v", err)
		return d
	}

	d.notificationSink = make(chan UdpNotification, 1024)
	d.connection = conn

	go d.readLoop()
	go d.writeLoop()

	return d
}

func (d *UdpRequestDispatcher) readLoop() {
	buf := make([]byte, 1024)
	for {
		n, addr, err := d.connection.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			return // graceful shutdown usually
		}
		packet := make([]byte, n)
		copy(packet, buf[:n])
		go d.HandleDatagramPackage(context.Background(), packet, addr)
	}
}

func (d *UdpRequestDispatcher) writeLoop() {
	for notification := range d.notificationSink {
		udpAddr, ok := notification.RecipientAddress.(*net.UDPAddr)
		if !ok {
			continue
		}
		data := GetBufferFromNotificationType(notification.Type)
		_, _ = d.connection.WriteToUDP(data, udpAddr)
	}
}

// @MappedFrom sendSignal(InetSocketAddress address, UdpNotificationType signal)
func (d *UdpRequestDispatcher) SendSignal(address net.Addr, signal UdpNotificationType) {
	if d.notificationSink != nil {
		notification := UdpNotification{
			RecipientAddress: address,
			Type:             signal,
		}
		select {
		case <-d.stopChan:
			return
		case d.notificationSink <- notification:
		default:
			// Fallback to goroutine to prevent dropping notifications when the buffer is full
			// This matches Java's unbounded sink behavior.
			go func() {
				select {
				case <-d.stopChan:
					return
				case d.notificationSink <- notification:
				}
			}()
		}
	}
}

// @MappedFrom handleDatagramPackage(DatagramPacket packet)
func (d *UdpRequestDispatcher) HandleDatagramPackage(ctx context.Context, packet []byte, senderAddress net.Addr) {
	req := d.ParseRequest(packet)
	if req == nil {
		d.sendErrorCode(senderAddress, constant.ResponseStatusCode_ILLEGAL_ARGUMENT)
		return
	}

	s := d.sessionService.GetLocalUserSession(ctx, req.UserID, req.DeviceType)
	if s == nil || s.ID != req.SessionID {
		d.sendErrorCode(senderAddress, constant.ResponseStatusCode_UPDATE_HEARTBEAT_OF_NONEXISTENT_SESSION)
		return
	}

	switch req.Type {
	case HeartbeatRequest:
		s.SetLastHeartbeatRequestTimestampToNow()
		// Java doesn't send a response for UDP heartbeat.
	case GoOfflineRequest:
		d.sessionService.UnregisterSession(ctx, req.UserID, req.DeviceType, nil, constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT)
		d.sendErrorCode(senderAddress, constant.ResponseStatusCode_OK)
	}
}

func (d *UdpRequestDispatcher) sendErrorCode(addr net.Addr, code constant.ResponseStatusCode) {
	if d.connection == nil {
		return
	}
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return
	}
	data := GetBufferFromStatusCode(code)
	_, _ = d.connection.WriteToUDP(data, udpAddr)
}

// @MappedFrom parseRequest(ByteBuf byteBuf)
func (d *UdpRequestDispatcher) ParseRequest(buffer []byte) *UdpSignalRequest {
	if len(buffer) < 13 { // 1 (type) + 8 (userId) + 1 (device) + 3 (only first 3 bytes of sessionId??)
		// Wait, Java sessionID is int, which is 4 bytes.
		// Go: binary.BigEndian.Uint32(buffer[10:14])
		return nil
	}
	if len(buffer) < 14 {
		return nil
	}
	reqType, err := ParseUdpRequestType(int(buffer[0]))
	if err != nil {
		return nil
	}
	userID := int64(binary.BigEndian.Uint64(buffer[1:9]))
	deviceType := protocol.DeviceType(buffer[9])
	sessionID := int64(binary.BigEndian.Uint32(buffer[10:14]))

	return NewUdpSignalRequest(reqType, userID, deviceType, sessionID)
}
