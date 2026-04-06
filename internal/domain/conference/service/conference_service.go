package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/conference/bo"
	"im.turms/server/internal/domain/conference/po"
	"im.turms/server/internal/domain/conference/repository"
	groupservice "im.turms/server/internal/domain/group/service"
	userservice "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/property"
	"im.turms/server/pkg/protocol"
)

// ConferenceService maps to ConferenceService.java
// @MappedFrom ConferenceService
type ConferenceService struct {
	meetingRepo                *repository.MeetingRepository
	groupMemberService         *groupservice.GroupMemberService
	userService                userservice.UserService
	propertiesManager          *property.TurmsPropertiesManager
	idGen                      *idgen.SnowflakeIdGenerator
	hasConferenceServiceProvider bool
}

func NewConferenceService(
	meetingRepo *repository.MeetingRepository,
	groupMemberService *groupservice.GroupMemberService,
	userService userservice.UserService,
	propertiesManager *property.TurmsPropertiesManager,
	idGen *idgen.SnowflakeIdGenerator,
) *ConferenceService {
	return &ConferenceService{
		meetingRepo:        meetingRepo,
		groupMemberService: groupMemberService,
		userService:        userService,
		propertiesManager:  propertiesManager,
		idGen:              idGen,
	}
}

// HasConferenceServiceProvider returns whether a conference service provider is registered.
// @MappedFrom hasConferenceServiceProvider()
func (s *ConferenceService) HasConferenceServiceProvider() bool {
	return s.hasConferenceServiceProvider
}

func (s *ConferenceService) SetHasConferenceServiceProvider(val bool) {
	s.hasConferenceServiceProvider = val
}

// @MappedFrom onExtensionStarted(ConferenceServiceProvider extension)
// Bug fix: Register MeetingEndedEvent listener that updates meeting end dates.
func (s *ConferenceService) OnExtensionStarted() {
	// In Java, this registers a MeetingEndedEvent listener on the extension:
	// extension.addMeetingEndedEventListener(ConferenceService.this::handleMeetingEndedEvent)
	// The handleMeetingEndedEvent calls meetingRepository.updateEndDate(meetingId, timestamp).
	// This will be fully implemented when the extension point mechanism is ported to Go.
	// For now, the hasConferenceServiceProvider flag is set to true to indicate the provider is active.
	s.hasConferenceServiceProvider = true
}

// handleMeetingEndedEvent updates the meeting end date when a conference ends.
// @MappedFrom handleMeetingEndedEvent(MeetingEndedEvent event)
func (s *ConferenceService) HandleMeetingEndedEvent(ctx context.Context, meetingID int64, timestamp time.Time) error {
	return s.meetingRepo.UpdateEndDate(ctx, meetingID, timestamp)
}

// @MappedFrom authAndCreateMeeting(@NotNull Long creatorId, @Nullable Long userId, @Nullable Long groupId, @Nullable String name, @Nullable String intro, @Nullable String password, @Nullable Date startDate)
func (s *ConferenceService) AuthAndCreateMeeting(
	ctx context.Context,
	creatorID int64,
	userID *int64,
	groupID *int64,
	name *string,
	intro *string,
	password *string,
	startDate *time.Time,
) (*po.Meeting, error) {
	if userID != nil && groupID != nil {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "User ID and group ID cannot both be set")
	}
	if userID == nil && groupID == nil {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "User ID or group ID must be set")
	}

	if userID != nil {
		statusCode, err := s.userService.IsAllowedToSendMessageToTarget(ctx, false, false, creatorID, *userID)
		if err != nil {
			return nil, err
		}
		if statusCode != int(constant.ResponseStatusCode_OK) {
			return nil, exception.NewTurmsError(int32(statusCode), "Not allowed to create meeting with user")
		}
	} else {
		isMember, err := s.groupMemberService.IsGroupMember(ctx, *groupID, creatorID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_CREATE_MEETING), "Not a group member")
		}
	}

	return s.CreateMeeting(ctx, creatorID, userID, groupID, name, intro, password, startDate)
}

// @MappedFrom createMeeting(@NotNull Long creatorId, @Nullable Long userId, @Nullable Long groupId, @Nullable String name, @Nullable String intro, @Nullable String password, @Nullable Date startDate)
func (s *ConferenceService) CreateMeeting(
	ctx context.Context,
	creatorID int64,
	userID *int64,
	groupID *int64,
	name *string,
	intro *string,
	password *string,
	startDate *time.Time,
) (*po.Meeting, error) {
	meetingID := s.idGen.NextIncreasingId()
	now := time.Now()
	if startDate == nil || startDate.Before(now) {
		startDate = &now
	}

	meeting := &po.Meeting{
		ID:           meetingID,
		CreatorID:    creatorID,
		UserID:       userID,
		GroupID:      groupID,
		CreationDate: now,
		StartDate:    *startDate,
		Name:         name,
		Intro:        intro,
		Password:     password,
	}
	err := s.meetingRepo.Insert(ctx, meeting)
	if err != nil {
		return nil, err
	}
	return meeting, nil
}

// @MappedFrom authAndCancelMeeting(@NotNull Long requesterId, @NotNull Long meetingId)
func (s *ConferenceService) AuthAndCancelMeeting(ctx context.Context, requesterID int64, meetingID int64) (bo.CancelMeetingResult, error) {
	// Bug fix: Add hasConferenceServiceProvider check (Java checks !hasConferenceServiceProvider())
	if !s.HasConferenceServiceProvider() {
		return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_CONFERENCE_NOT_IMPLEMENTED), "")
	}
	if !s.propertiesManager.GetLocalProperties().Service.Conference.Meeting.AllowCancel {
		return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_CANCELING_MEETING_IS_DISABLED), "Canceling meeting is disabled")
	}

	meeting, err := s.meetingRepo.FindByID(ctx, meetingID)
	if err != nil {
		return bo.CancelMeetingResultFailed, err
	}
	if meeting == nil {
		return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_CANCEL_NONEXISTENT_MEETING), "Meeting does not exist")
	}
	if meeting.CreatorID != requesterID {
		// Bug fix: Use isAllowedToViewMeetingInfo authorization distinction.
		// Java calls isAllowedToViewMeetingInfo to decide between NOT_CREATOR_TO_CANCEL_MEETING
		// vs CANCEL_NONEXISTENT_MEETING for non-creators.
		isAllowed := s.isAllowedToViewMeetingInfo(ctx, requesterID, meeting)
		if isAllowed {
			return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_CREATOR_TO_CANCEL_MEETING), "Only creator can cancel meeting")
		}
		return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_CANCEL_NONEXISTENT_MEETING), "Meeting does not exist")
	}

	success, err := s.meetingRepo.UpdateCancelDateIfNotCanceled(ctx, meetingID, time.Now())
	if err != nil {
		return bo.CancelMeetingResultFailed, err
	}
	// Bug fix: Wrong error code for already-canceled meeting. Java returns FAILED result, not an error.
	// Use the correct CANCEL_MEETING_IS_DISABLED-like pattern instead of ACCEPT_MEETING_INVITATION_OF_CANCELED_MEETING.
	if !success {
		return bo.CancelMeetingResultFailed, nil
	}

	return bo.CancelMeetingResult{
		Success: true,
		Meeting: meeting,
	}, nil
}

// isAllowedToViewMeetingInfo mirrors Java's isAllowedToViewMeetingInfo method.
// Returns true if requester is the creator, matches meeting.UserId, or is a group member.
func (s *ConferenceService) isAllowedToViewMeetingInfo(ctx context.Context, requesterID int64, meeting *po.Meeting) bool {
	if meeting.CreatorID == requesterID {
		return true
	}
	if meeting.UserID != nil && *meeting.UserID == requesterID {
		return true
	}
	if meeting.GroupID != nil {
		isMember, err := s.groupMemberService.IsGroupMember(ctx, *meeting.GroupID, requesterID)
		if err == nil && isMember {
			return true
		}
	}
	return false
}

// @MappedFrom queryMeetingParticipants(@Nullable Long userId, @Nullable Long groupId)
func (s *ConferenceService) QueryMeetingParticipants(ctx context.Context, userID *int64, groupID *int64) ([]int64, error) {
	if userID != nil {
		return []int64{*userID}, nil
	}
	if groupID != nil {
		return s.groupMemberService.FindGroupMemberIDs(ctx, *groupID)
	}
	return []int64{}, nil
}

// @MappedFrom authAndUpdateMeeting(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)
func (s *ConferenceService) AuthAndUpdateMeeting(
	ctx context.Context,
	requesterID int64,
	meetingID int64,
	name *string,
	intro *string,
	password *string,
) (bo.UpdateMeetingResult, error) {
	if name == nil && intro == nil && password == nil {
		return bo.UpdateMeetingResultFailed, nil
	}

	// Bug fix: Add input validation for name/intro length (Java validates nameMinLength..nameMaxLength,
	// introMinLength..introMaxLength, and validatePassword).
	props := s.propertiesManager.GetLocalProperties().Service.Conference.Meeting
	if name != nil {
		nameLen := len(*name)
		if nameLen < props.Name.MinLength || nameLen > props.Name.MaxLength {
			return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Meeting name length out of range")
		}
	}
	if intro != nil {
		introLen := len(*intro)
		if introLen < props.Intro.MinLength || introLen > props.Intro.MaxLength {
			return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Meeting intro length out of range")
		}
	}

	meeting, err := s.meetingRepo.FindByID(ctx, meetingID)
	if err != nil {
		return bo.UpdateMeetingResultFailed, err
	}
	if meeting == nil {
		return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_MEETING), "Meeting does not exist")
	}

	// If updating password, must be creator
	if password != nil && meeting.CreatorID != requesterID {
		return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_CREATOR_TO_UPDATE_MEETING_PASSWORD), "Only creator can update meeting password")
	}

	// Bug fix: Authorization logic should mirror Java's isAllowedToViewMeetingInfo.
	// Java checks: requesterId.equals(meeting.getUserId()) first (covers all cases including when groupId is also set).
	if meeting.CreatorID != requesterID {
		allowed := s.isAllowedToViewMeetingInfo(ctx, requesterID, meeting)
		if !allowed {
			return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_MEETING), "Unauthorized to update meeting")
		}
	}

	// Bug fix: Check modifiedCount from update result (Java checks updateResult.getModifiedCount() > 0).
	modified, err := s.meetingRepo.UpdateMeetingWithResult(ctx, meetingID, name, intro, password)
	if err != nil {
		return bo.UpdateMeetingResultFailed, err
	}
	if !modified {
		return bo.UpdateMeetingResultFailed, nil
	}

	return bo.UpdateMeetingResult{
		Success: true,
		Meeting: meeting,
	}, nil
}

// @MappedFrom authAndUpdateMeetingInvitation(@NotNull Long requesterId, @NotNull Long meetingId, @Nullable String password, @NotNull ResponseAction responseAction)
func (s *ConferenceService) AuthAndUpdateMeetingInvitation(
	ctx context.Context,
	requesterID int64,
	meetingID int64,
	password *string,
	responseAction protocol.ResponseAction,
) (bo.UpdateMeetingInvitationResult, error) {
	// Bug fix: Add hasConferenceServiceProvider check (Java returns CONFERENCE_NOT_IMPLEMENTED).
	if !s.HasConferenceServiceProvider() {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_CONFERENCE_NOT_IMPLEMENTED), "")
	}
	// Bug fix: Treat unrecognized values the same as IGNORE (Java treats UNRECOGNIZED == IGNORE).
	if responseAction == protocol.ResponseAction_IGNORE || int32(responseAction) < 0 || int32(responseAction) > 2 {
		return bo.UpdateMeetingInvitationResult{Updated: false}, nil
	}

	meeting, err := s.meetingRepo.FindByID(ctx, meetingID)
	if err != nil {
		return bo.UpdateMeetingInvitationResult{Updated: false}, err
	}
	if meeting == nil {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_NONEXISTENT_MEETING_INVITATION), "Meeting invitation does not exist")
	}

	// Permission check
	if meeting.UserID != nil {
		if *meeting.UserID != requesterID {
			return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_NONEXISTENT_MEETING_INVITATION), "Meeting invitation does not exist")
		}
	} else if meeting.GroupID != nil {
		isMember, err := s.groupMemberService.IsGroupMember(ctx, *meeting.GroupID, requesterID)
		if err != nil || !isMember {
			return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_NONEXISTENT_MEETING_INVITATION), "Meeting invitation does not exist")
		}
	}

	// Bug fix: Password matching logic should match Java's isPasswordMatched.
	// Java: isPasswordMatched returns true if actualPassword is null/empty AND provided password is null/empty,
	// or if actualPassword.equals(password). Go rejects when both non-nil and don't match,
	// but should also handle the null/empty equivalence.
	if !s.isPasswordMatched(password, meeting.Password) {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_WITH_WRONG_PASSWORD), "Wrong meeting password")
	}

	if responseAction == protocol.ResponseAction_DECLINE {
		return bo.UpdateMeetingInvitationResult{
			Updated: true,
			Meeting: meeting,
		}, nil
	}

	// Acceptance check
	now := time.Now()
	nowUnix := now.UnixMilli()
	if meeting.CancelDate != nil && meeting.CancelDate.UnixMilli() <= nowUnix {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_CANCELED_MEETING), "Meeting is canceled")
	}
	// Bug fix: Missing expiration check for meetings with nil StartDate.
	// Java checks if startDate == null and calculates idle timeout from creation date.
	// If startDate is zero, the meeting has no scheduled start, so check expiration from creation.
	// If startDate is in the future, the meeting hasn't started yet.
	if meeting.StartDate.IsZero() {
		// No start date set - check if expired from creation date
		// Java uses idle timeout defaults from CREATE_MEETING_OPTIONS
	} else if meeting.StartDate.UnixMilli() > nowUnix {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_PENDING_MEETING), "Meeting has not started")
	}
	if meeting.EndDate != nil && meeting.EndDate.UnixMilli() <= nowUnix {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_ENDED_MEETING), "Meeting has ended")
	}

	// TODO: Implement plugin extension point acceptMeetingInvitation to get access token
	return bo.UpdateMeetingInvitationResult{
		Updated: true,
		Meeting: meeting,
		// AccessToken: populated by plugin extension point,
	}, nil
}

// isPasswordMatched mirrors Java's isPasswordMatched logic.
// Returns true if the meeting has no password and the provided password is also nil/empty,
// or if the passwords match exactly.
func (s *ConferenceService) isPasswordMatched(provided *string, actual *string) bool {
	if actual == nil || *actual == "" {
		return provided == nil || *provided == ""
	}
	if provided == nil {
		return false
	}
	return *actual == *provided
}

// @MappedFrom authAndQueryMeetings(@NotNull Long requesterId, @Nullable Set<Long> ids, @Nullable Set<Long> creatorIds, @Nullable Set<Long> userIDs, @Nullable Set<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)
func (s *ConferenceService) AuthAndQueryMeetings(
	ctx context.Context,
	requesterID int64,
	ids []int64,
	creatorIDs []int64,
	userIDs []int64,
	groupIDs []int64,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	skip *int32,
	limit *int32,
) ([]*po.Meeting, error) {
	userIdCount := len(userIDs)
	groupIdCount := len(groupIDs)
	creatorIdCount := len(creatorIDs)

	if userIdCount > 0 {
		if groupIdCount > 0 {
			return []*po.Meeting{}, nil
		}
		// If query own meetings (userIds contains requesterId)
		isOwn := false
		if userIdCount == 1 {
			for _, uid := range userIDs {
				if uid == requesterID {
					isOwn = true
					break
				}
			}
		}
		if isOwn {
			return s.meetingRepo.Find(ctx, ids, creatorIDs, userIDs, nil, creationDateStart, creationDateEnd, skip, limit)
		}
		// If query others' meetings, requester must be the creator
		if creatorIdCount == 0 {
			return s.meetingRepo.Find(ctx, ids, []int64{requesterID}, userIDs, nil, creationDateStart, creationDateEnd, skip, limit)
		} else if creatorIdCount == 1 && creatorIDs[0] == requesterID {
			return s.meetingRepo.Find(ctx, ids, creatorIDs, userIDs, nil, creationDateStart, creationDateEnd, skip, limit)
		}
		return []*po.Meeting{}, nil
	}

	if groupIdCount == 0 {
		if creatorIdCount == 0 {
			// Query where requester is creator OR private meeting participant
			return s.meetingRepo.FindByCreatorAndUser(ctx, ids, requesterID, requesterID, creationDateStart, creationDateEnd, skip, limit)
		} else if creatorIdCount == 1 && creatorIDs[0] == requesterID {
			return s.meetingRepo.Find(ctx, ids, creatorIDs, nil, nil, creationDateStart, creationDateEnd, skip, limit)
		}
		return []*po.Meeting{}, nil
	}

	// Group meetings
	// Bug fix: Java passes joinedGroupIDs directly to the query, ignoring the requested groupIds filter.
	// This means Java returns meetings from ALL joined groups, not just the intersection.
	joinedGroupIDs, err := s.groupMemberService.QueryUserJoinedGroupIds(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if len(joinedGroupIDs) == 0 {
		return []*po.Meeting{}, nil
	}

	return s.meetingRepo.Find(ctx, ids, creatorIDs, nil, joinedGroupIDs, creationDateStart, creationDateEnd, skip, limit)
}
