package proto

import "im.turms/server/pkg/protocol"

// SimpleTurmsNotification maps to SimpleTurmsNotification in Java.
// @MappedFrom SimpleTurmsNotification
type SimpleTurmsNotification struct {
	RequesterID       int64
	CloseStatus       *int32
	RelayedRequestType *protocol.TurmsRequest_Kind
}

func NewSimpleTurmsNotification(requesterID int64, closeStatus *int32, relayedRequestType *protocol.TurmsRequest_Kind) *SimpleTurmsNotification {
	return &SimpleTurmsNotification{
		RequesterID:       requesterID,
		CloseStatus:       closeStatus,
		RelayedRequestType: relayedRequestType,
	}
}

// SimpleTurmsRequest maps to SimpleTurmsRequest in Java.
// @MappedFrom SimpleTurmsRequest
type SimpleTurmsRequest struct {
	RequestID            int64
	Type                 *protocol.TurmsRequest_Kind
	CreateSessionRequest *protocol.CreateSessionRequest
}

func NewSimpleTurmsRequest(requestID int64, reqType *protocol.TurmsRequest_Kind, createSessionReq *protocol.CreateSessionRequest) *SimpleTurmsRequest {
	return &SimpleTurmsRequest{
		RequestID:            requestID,
		Type:                 reqType,
		CreateSessionRequest: createSessionReq,
	}
}

// @MappedFrom toString()
func (r *SimpleTurmsRequest) ToString() string {
	return ""
}

// TurmsNotificationParser maps to TurmsNotificationParser in Java.
// @MappedFrom TurmsNotificationParser
type TurmsNotificationParser struct{}

// @MappedFrom parseSimpleNotification(CodedInputStream turmsRequestInputStream)
func (p *TurmsNotificationParser) ParseSimpleNotification(turmsRequestInputStream []byte) (*SimpleTurmsNotification, error) {
	// Stub implementation
	return nil, nil
}

// TurmsRequestParser maps to TurmsRequestParser in Java.
// @MappedFrom TurmsRequestParser
type TurmsRequestParser struct{}

// @MappedFrom parseSimpleRequest(CodedInputStream turmsRequestInputStream)
func (p *TurmsRequestParser) ParseSimpleRequest(turmsRequestInputStream []byte) (*SimpleTurmsRequest, error) {
	// Stub implementation
	return nil, nil
}
