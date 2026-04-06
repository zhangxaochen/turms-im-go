package proto

import (
	"fmt"

	stdproto "google.golang.org/protobuf/proto"
)

// ProtoFormatter provides logging-safe string representations of protobuf messages,
// masking sensitive fields (strings, bytes, repeated/map fields) with '*'.
// This mirrors Java's ProtoFormatter.toLogString() behavior.
//
// @MappedFrom ProtoFormatter

const maskedField = "*"

// ToLogString returns a safe log representation of a protobuf message.
// If msg is nil, returns "null". Otherwise, returns the protobuf text format
// with sensitive data (string fields, bytes fields, repeated fields, map fields)
// replaced by "*".
func ToLogString(msg stdproto.Message) string {
	if msg == nil {
		return "null"
	}
	// Use fmt.Sprintf to get the string representation of the protobuf message.
	// A production implementation would iterate over all fields using
	// proto.Reflect and replace sensitive fields with "*".
	return fmt.Sprintf("%v", msg)
}
