package util

import "im.turms/server/internal/domain/user/po"

// ExpirableRequestInspector maps to ExpirableRequestInspector.java
// @MappedFrom ExpirableRequestInspector
type ExpirableRequestInspector struct {
}

// @MappedFrom isProcessedByResponder(@Nullable RequestStatus status)
func (i *ExpirableRequestInspector) IsProcessedByResponder(status po.RequestStatus) bool {
	return status == po.RequestStatusAccepted || status == po.RequestStatusDeclined || status == po.RequestStatusIgnored
}
