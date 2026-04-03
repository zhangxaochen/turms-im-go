package repository

// MeetingRepository maps to MeetingRepository.java
// @MappedFrom MeetingRepository
type MeetingRepository struct {
}

// @MappedFrom updateEndDate(Long meetingId, Date endDate)
func (r *MeetingRepository) UpdateEndDate() {
}

// @MappedFrom updateCancelDateIfNotCanceled(Long meetingId, Date cancelDate)
func (r *MeetingRepository) UpdateCancelDateIfNotCanceled() {
}

// @MappedFrom updateMeeting(Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)
func (r *MeetingRepository) UpdateMeeting() {
}

// @MappedFrom find(@Nullable Collection<Long> ids, @Nullable Collection<Long> creatorIds, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)
func (r *MeetingRepository) Find() {
}

// @MappedFrom find(@Nullable Collection<Long> ids, @NotNull Long creatorId, @NotNull Long userId, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)
func (r *MeetingRepository) FindByCreatorAndUser() {
}
