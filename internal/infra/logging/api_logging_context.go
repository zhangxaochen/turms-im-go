package logging

import "math/rand/v2"

// LoggingRequestProperties holds per-request-type logging configuration.
// @MappedFrom im.turms.server.common.infra.property.env.service.env.clientapi.property.LoggingRequestProperties
type LoggingRequestProperties struct {
	SampleRate float32
}

// ApiLoggingContext maps to ApiLoggingContext in Java (turms-gateway).
// @MappedFrom ApiLoggingContext
type ApiLoggingContext struct {
	// typeToSupportedLoggingRequestProperties maps request kind (int) to its logging properties.
	// An absent key means the request type is excluded/not logged.
	typeToSupportedLoggingRequestProperties map[int]*LoggingRequestProperties

	// typeToSupportedLoggingNotificationsProperties maps notification kind to its logging properties.
	typeToSupportedLoggingNotificationsProperties map[int]*LoggingRequestProperties

	// heartbeatSampleRate is the fraction [0,1] of heartbeat requests to log.
	heartbeatSampleRate float32
}

func NewApiLoggingContext(
	requestProps map[int]*LoggingRequestProperties,
	notificationProps map[int]*LoggingRequestProperties,
	heartbeatSampleRate float32,
) *ApiLoggingContext {
	return &ApiLoggingContext{
		typeToSupportedLoggingRequestProperties:       requestProps,
		typeToSupportedLoggingNotificationsProperties: notificationProps,
		heartbeatSampleRate:                           heartbeatSampleRate,
	}
}

// NewApiLoggingContextDefault creates a stub instance with no filtering (logs everything).
// Use NewApiLoggingContext to configure properly.
func NewApiLoggingContextDefault() *ApiLoggingContext {
	return &ApiLoggingContext{
		typeToSupportedLoggingRequestProperties:       nil,
		typeToSupportedLoggingNotificationsProperties: nil,
		heartbeatSampleRate:                           1.0,
	}
}

// shouldLogBySampleRate implements the probabilistic sampling logic from BaseApiLoggingContext.shouldLog(float).
// Returns false if rate <= 0, true if rate >= 1, otherwise random comparison.
func shouldLogBySampleRate(sampleRate float32) bool {
	if sampleRate <= 0 {
		return false
	}
	if sampleRate >= 1.0 {
		return true
	}
	return rand.Float32() < sampleRate
}

// shouldLogForType checks the map for requestType and applies its sample rate.
// Returns false if the type is not in the map.
func shouldLogForType(requestType int, props map[int]*LoggingRequestProperties) bool {
	if props == nil {
		// nil map means log nothing (type not configured)
		return false
	}
	p, ok := props[requestType]
	if !ok {
		return false
	}
	return shouldLogBySampleRate(p.SampleRate)
}

// @MappedFrom shouldLogHeartbeatRequest()
func (c *ApiLoggingContext) ShouldLogHeartbeatRequest() bool {
	return shouldLogBySampleRate(c.heartbeatSampleRate)
}

// @MappedFrom shouldLogRequest(TurmsRequest.KindCase requestType)
func (c *ApiLoggingContext) ShouldLogRequest(requestType int) bool {
	return shouldLogForType(requestType, c.typeToSupportedLoggingRequestProperties)
}

// @MappedFrom shouldLogNotification(TurmsRequest.KindCase requestType)
func (c *ApiLoggingContext) ShouldLogNotification(requestType int) bool {
	return shouldLogForType(requestType, c.typeToSupportedLoggingNotificationsProperties)
}
