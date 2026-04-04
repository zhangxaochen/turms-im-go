package bo

import (
	"encoding/json"
	"net"
	"time"
)
type BlockedClient[T any] struct {
	ID                 T
	BlockEndTimeMillis int64
}

// MarshalJSON implements custom JSON serialization equivalent to BlockedClientSerializer.java
// @MappedFrom serialize(BlockedClient value, JsonGenerator gen, SerializerProvider provider)
func (b BlockedClient[T]) MarshalJSON() ([]byte, error) {
	// Format time
	t := time.UnixMilli(b.BlockEndTimeMillis).UTC()
	timeStr := t.Format("2006-01-02 15:04:05.000")
	var idVal interface{} = b.ID
	
	// Try to handle IP bytes to string if it was somehow storing []byte, but type T is usually string or int64 here
	if bytesVal, ok := any(b.ID).([]byte); ok {
		idVal = net.IP(bytesVal).String()
	}

	return json.Marshal(map[string]interface{}{
		"_id":          idVal,
		"blockEndTime": timeStr,
	})
}
