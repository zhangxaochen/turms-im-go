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
	meetingRepo        *repository.MeetingRepository
	groupMemberService *groupservice.GroupMemberService
	userService        userservice.UserService
	propertiesManager  *property.TurmsPropertiesManager
	idGen              *idgen.SnowflakeIdGenerator
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

// @MappedFrom onExtensionStarted(ConferenceServiceProvider extension)
func (s *ConferenceService) OnExtensionStarted() {
	// Plugin logic can be added here once the extension point mechanism is fully implemented in Go.
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
		return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_CREATOR_TO_CANCEL_MEETING), "Only creator can cancel meeting")
	}

	success, err := s.meetingRepo.UpdateCancelDateIfNotCanceled(ctx, meetingID, time.Now())
	if err != nil {
		return bo.CancelMeetingResultFailed, err
	}
	if !success {
		return bo.CancelMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_CANCELED_MEETING), "Meeting is already canceled")
	}

	return bo.CancelMeetingResult{
		Success: true,
		Meeting: meeting,
	}, nil
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

	// Permission check
	if meeting.CreatorID != requesterID {
		if meeting.GroupID != nil {
			isMember, err := s.groupMemberService.IsGroupMember(ctx, *meeting.GroupID, requesterID)
			if err != nil || !isMember {
				return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_MEETING), "Unauthorized to update meeting")
			}
		} else if meeting.UserID != nil && *meeting.UserID != requesterID {
			return bo.UpdateMeetingResultFailed, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_MEETING), "Unauthorized to update meeting")
		}
	}

	err = s.meetingRepo.UpdateMeeting(ctx, meetingID, name, intro, password)
	if err != nil {
		return bo.UpdateMeetingResultFailed, err
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
	if responseAction == protocol.ResponseAction_IGNORE {
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

	if password != nil && (meeting.Password != nil && *password != *meeting.Password) {
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
	if meeting.StartDate.UnixMilli() > nowUnix {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_PENDING_MEETING), "Meeting has not started")
	}
	if meeting.EndDate != nil && meeting.EndDate.UnixMilli() <= nowUnix {
		return bo.UpdateMeetingInvitationResult{Updated: false}, exception.NewTurmsError(int32(constant.ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_ENDED_MEETING), "Meeting has ended")
	}

	// TODO: Implement LiveKit or other conference provider integration here to get access token
	return bo.UpdateMeetingInvitationResult{
		Updated: true,
		Meeting: meeting,
		// AccessToken: tokenFromProvider,
	}, nil
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
	joinedGroupIDs, err := s.groupMemberService.QueryUserJoinedGroupIds(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if len(joinedGroupIDs) == 0 {
		return []*po.Meeting{}, nil
	}

	joinedMap := make(map[int64]bool)
	for _, gid := range joinedGroupIDs {
		joinedMap[gid] = true
	}
	validGroupIDs := make([]int64, 0)
	for _, gid := range groupIDs {
		if joinedMap[gid] {
			validGroupIDs = append(validGroupIDs, gid)
		}
	}
	if len(validGroupIDs) == 0 {
		return []*po.Meeting{}, nil
	}

	return s.meetingRepo.Find(ctx, ids, creatorIDs, nil, validGroupIDs, creationDateStart, creationDateEnd, skip, limit)
}
