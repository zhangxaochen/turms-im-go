package controller

import (
	"im.turms/server/pkg/protocol"
	"google.golang.org/protobuf/proto"
)

func buildSuccessNotification(reqID *int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		RequestId: reqID,
		Code:      proto.Int32(1000), // SUCCESS
	}
}
