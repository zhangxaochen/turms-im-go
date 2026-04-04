package bo

import (
	"im.turms/server/internal/domain/conference/po"
)

// @MappedFrom CancelMeetingResult(boolean success, @Nullable Meeting meeting)
type CancelMeetingResult struct {
	Success bool
	Meeting *po.Meeting
}

var CancelMeetingResultFailed = CancelMeetingResult{Success: false, Meeting: nil}

// @MappedFrom UpdateMeetingInvitationResult(boolean updated, @Nullable String accessToken, @Nullable Meeting meeting)
type UpdateMeetingInvitationResult struct {
	Updated     bool
	AccessToken *string
	Meeting     *po.Meeting
}

// @MappedFrom UpdateMeetingResult(boolean success, @Nullable Meeting meeting)
type UpdateMeetingResult struct {
	Success bool
	Meeting *po.Meeting
}

var UpdateMeetingResultFailed = UpdateMeetingResult{Success: false, Meeting: nil}
