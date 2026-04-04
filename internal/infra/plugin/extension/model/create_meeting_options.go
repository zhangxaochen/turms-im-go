package model

// CreateMeetingOptions represents the options for creating a meeting.
// @MappedFrom CreateMeetingOptions.java
type CreateMeetingOptions struct {
	MaxParticipants   *int32
	IdleTimeoutMillis *int64
}
