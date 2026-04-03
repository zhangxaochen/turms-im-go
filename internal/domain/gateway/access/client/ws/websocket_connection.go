package ws

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
)

// @MappedFrom WebSocketConnection
type WebSocketConnection struct {
	*common.BaseNetConnection
	conn         *websocket.Conn
	closeTimeout time.Duration
}

func NewWebSocketConnection(conn *websocket.Conn, isConnected bool, closeTimeout time.Duration) *WebSocketConnection {
	return &WebSocketConnection{
		BaseNetConnection: common.NewBaseNetConnection(isConnected),
		conn:              conn,
		closeTimeout:      closeTimeout,
	}
}

// @MappedFrom getAddress()
func (c *WebSocketConnection) GetAddress() net.Addr {
	return c.conn.RemoteAddr()
}

// @MappedFrom send(ByteBuf buffer)
func (c *WebSocketConnection) Send(ctx context.Context, buffer []byte) error {
	err := c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, buffer)
}

// @MappedFrom close(CloseReason closeReason)
func (c *WebSocketConnection) CloseWithReason(reason common.CloseReason) error {
	if !c.IsConnected() {
		return nil
	}
	c.BaseNetConnection.CloseWithReason(reason)

	if reason.Status != constant.SessionCloseStatus_UNKNOWN_ERROR {
		go func() {
			for i := 0; i < 3; i++ {
				// Gorilla websocket allows sending a close framing message
				// Here we send normal closure to gracefully stop
				err := c.conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason.Reason),
					time.Now().Add(time.Second),
				)
				if err == nil {
					break
				}
				log.Printf("Failed to send WS close notification attempt %d: %v", i+1, err)
				time.Sleep(3 * time.Second)
			}

			if c.closeTimeout > 0 {
				time.Sleep(c.closeTimeout)
			}

			err := c.conn.Close()
			if err != nil {
				log.Printf("Failed to close the WS connection %s: %v", c.GetAddress(), err)
			}
		}()
		return nil
	}

	err := c.conn.Close()
	if err != nil {
		log.Printf("Failed to close the WS connection %s: %v", c.GetAddress(), err)
	}
	return err
}

// @MappedFrom close()
func (c *WebSocketConnection) Close() error {
	if !c.IsConnected() {
		return nil
	}
	c.BaseNetConnection.Close()
	err := c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second),
	)
	if err != nil {
		log.Printf("Failed to write WS close message to %s: %v", c.GetAddress(), err)
	}
	
	err = c.conn.Close()
	if err != nil {
		log.Printf("Failed to close the WS connection %s: %v", c.GetAddress(), err)
	}
	return err
}
