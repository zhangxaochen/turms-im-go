package bo

// @MappedFrom CancelMeetingResult(boolean success, @Nullable Meeting meeting)
type CancelMeetingResult struct {
	Success bool
	Meeting interface{} // Replace with actual Meeting type
}

// @MappedFrom UpdateMeetingInvitationResult(boolean updated, @Nullable String accessToken, @Nullable Meeting meeting)
type UpdateMeetingInvitationResult struct {
	Updated     bool
	AccessToken *string
	Meeting     interface{}
}

// @MappedFrom UpdateMeetingResult(boolean success, @Nullable Meeting meeting)
type UpdateMeetingResult struct {
	Success bool
	Meeting interface{}
}
