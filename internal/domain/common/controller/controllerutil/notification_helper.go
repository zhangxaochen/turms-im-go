package controllerutil

import (
	"google.golang.org/protobuf/proto"

	"im.turms/server/pkg/protocol"
)

// BuildSuccessNotification constructs a minimal success TurmsNotification (code 1000).
func BuildSuccessNotification(reqID *int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		RequestId: reqID,
		Code:      proto.Int32(1000),
	}
}

// BuildDataLongNotification constructs a TurmsNotification containing a single int64 data value.
func BuildDataLongNotification(reqID *int64, value int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		RequestId: reqID,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: value,
			},
		},
	}
}

// BuildDataLongsNotification constructs a TurmsNotification containing multiple int64 data values.
func BuildDataLongsNotification(reqID *int64, values []int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		RequestId: reqID,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
				LongsWithVersion: &protocol.LongsWithVersion{
					Longs: values,
				},
			},
		},
	}
}
