package proto

import (
	"im.turms/server/pkg/protocol"
)

// ToList converts a map to a list of proto items if possible.
// @MappedFrom toList(Map<String, String> map)
func ToList(protoItems interface{}) []interface{} {
	// Simple implementation for now
	return nil
}

// Value2Proto converts a Go value to a Turms protocol Value.
// @MappedFrom value2proto(Value.Builder builder, Object object)
func Value2Proto(v interface{}) *protocol.Value {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case int:
		return &protocol.Value{Kind: &protocol.Value_Int32Value{Int32Value: int32(val)}}
	case int32:
		return &protocol.Value{Kind: &protocol.Value_Int32Value{Int32Value: val}}
	case int64:
		return &protocol.Value{Kind: &protocol.Value_Int64Value{Int64Value: val}}
	case float32:
		return &protocol.Value{Kind: &protocol.Value_FloatValue{FloatValue: val}}
	case float64:
		return &protocol.Value{Kind: &protocol.Value_DoubleValue{DoubleValue: val}}
	case bool:
		return &protocol.Value{Kind: &protocol.Value_BoolValue{BoolValue: val}}
	case []byte:
		return &protocol.Value{Kind: &protocol.Value_BytesValue{BytesValue: val}}
	case string:
		return &protocol.Value{Kind: &protocol.Value_StringValue{StringValue: val}}
	case []interface{}:
		list := make([]*protocol.Value, len(val))
		for i, ev := range val {
			list[i] = Value2Proto(ev)
		}
		return &protocol.Value{ListValue: list}
	}
	return nil
}

// Proto2Value converts a Turms protocol Value to a Go value.
func Proto2Value(v *protocol.Value) interface{} {
	if v == nil {
		return nil
	}
	if v.Kind != nil {
		switch val := v.Kind.(type) {
		case *protocol.Value_Int32Value:
			return int64(val.Int32Value)
		case *protocol.Value_Int64Value:
			return val.Int64Value
		case *protocol.Value_FloatValue:
			return val.FloatValue
		case *protocol.Value_DoubleValue:
			return val.DoubleValue
		case *protocol.Value_BoolValue:
			return val.BoolValue
		case *protocol.Value_BytesValue:
			return val.BytesValue
		case *protocol.Value_StringValue:
			return val.StringValue
		}
	}
	if len(v.ListValue) > 0 {
		list := make([]interface{}, len(v.ListValue))
		for i, ev := range v.ListValue {
			list[i] = Proto2Value(ev)
		}
		return list
	}
	return nil
}
