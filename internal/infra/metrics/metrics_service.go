package metrics

import (
	"im.turms/server/pkg/protocol"
)

// MetricsService is used for recording metrics of handling requests.
type MetricsService interface {
	RecordRequest(requestType *protocol.TurmsRequest, size int, processingTimeMilli int64)
}

type metricsService struct{}

func NewMetricsService() MetricsService {
	return &metricsService{}
}

func (s *metricsService) RecordRequest(requestType *protocol.TurmsRequest, size int, processingTimeMilli int64) {
	// TODO: Integrate with prometheus metrics logger here, measuring counter/latency depending on request type.
}
