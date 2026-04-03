package service

// ConferenceService maps to ConferenceService.java
// @MappedFrom ConferenceService
type ConferenceService struct {
}

// @MappedFrom onExtensionStarted(ConferenceServiceProvider extension)
func (s *ConferenceService) OnExtensionStarted() {
}

// @MappedFrom authAndCancelMeeting(@NotNull Long requesterId, @NotNull Long meetingId)
func (s *ConferenceService) AuthAndCancelMeeting() {
}

// @MappedFrom queryMeetingParticipants(@Nullable Long userId, @Nullable Long groupId)
func (s *ConferenceService) QueryMeetingParticipants() {
}

// @MappedFrom authAndUpdateMeeting(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)
func (s *ConferenceService) AuthAndUpdateMeeting() {
}

// @MappedFrom authAndUpdateMeetingInvitation(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String password, @NotNull ResponseAction responseAction)
func (s *ConferenceService) AuthAndUpdateMeetingInvitation() {
}

// @MappedFrom authAndQueryMeetings(@NotNull Long requesterId, @Nullable Set<Long> ids, @Nullable Set<Long> creatorIds, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)
func (s *ConferenceService) AuthAndQueryMeetings() {
}
