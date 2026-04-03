package controller

import (
	"google.golang.org/protobuf/proto"
	"im.turms/server/pkg/protocol"
)

func buildSuccessNotification(reqID *int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		RequestId: reqID,
		Code:      proto.Int32(1000), // SUCCESS
	}
}
