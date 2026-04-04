package po

import (
	"time"

	"im.turms/server/pkg/protocol"
)

// Expirable maps to Expirable.java
// It represents a PO that can expire and holds a status.
// @MappedFrom Expirable
type Expirable interface {
	GetCreationDate() time.Time
	GetStatus() protocol.RequestStatus
	SetStatus(status protocol.RequestStatus)
}
