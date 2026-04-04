package controller

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/conference/service"
	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

// ConferenceServiceController maps to ConferenceServiceController.java
// @MappedFrom ConferenceServiceController
type ConferenceServiceController struct {
	conferenceService *service.ConferenceService
}

func NewConferenceServiceController(conferenceService *service.ConferenceService) *ConferenceServiceController {
	return &ConferenceServiceController{
		conferenceService: conferenceService,
	}
}

func (c *ConferenceServiceController) RegisterRoutes(r *router.Router) {
	r.RegisterController(&protocol.TurmsRequest_CreateMeetingRequest{}, c.HandleCreateMeetingRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteMeetingRequest{}, c.HandleDeleteMeetingRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateMeetingRequest{}, c.HandleUpdateMeetingRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryMeetingsRequest{}, c.HandleQueryMeetingsRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateMeetingInvitationRequest{}, c.HandleUpdateMeetingInvitationRequest)
}

// @MappedFrom handleCreateMeetingRequest(@NotNull ClientRequest<CreateMeetingRequest> clientRequest)
func (c *ConferenceServiceController) HandleCreateMeetingRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createMeetingRequest := req.GetCreateMeetingRequest()
	var startDate *time.Time
	if createMeetingRequest.StartDate != nil {
		t := time.UnixMilli(*createMeetingRequest.StartDate)
		startDate = &t
	}
	meeting, err := c.conferenceService.AuthAndCreateMeeting(
		ctx,
		s.UserID,
		createMeetingRequest.UserId,
		createMeetingRequest.GroupId,
		createMeetingRequest.Name,
		createMeetingRequest.Intro,
		createMeetingRequest.Password,
		startDate,
	)
	if err != nil {
		return nil, err
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: meeting.ID,
			},
		},
	}, nil
}

// @MappedFrom handleDeleteMeetingRequest(@NotNull ClientRequest<DeleteMeetingRequest> clientRequest)
func (c *ConferenceServiceController) HandleDeleteMeetingRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteMeetingRequest := req.GetDeleteMeetingRequest()
	result, err := c.conferenceService.AuthAndCancelMeeting(ctx, s.UserID, deleteMeetingRequest.GetId())
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, nil // Should be handled by error
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
	}, nil
}

// @MappedFrom handleQueryMeetingsRequest(@NotNull ClientRequest<QueryMeetingsRequest> clientRequest)
func (c *ConferenceServiceController) HandleQueryMeetingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryMeetingsRequest := req.GetQueryMeetingsRequest()
	var creationDateStart *time.Time
	if queryMeetingsRequest.CreationDateStart != nil {
		t := time.UnixMilli(*queryMeetingsRequest.CreationDateStart)
		creationDateStart = &t
	}
	var creationDateEnd *time.Time
	if queryMeetingsRequest.CreationDateEnd != nil {
		t := time.UnixMilli(*queryMeetingsRequest.CreationDateEnd)
		creationDateEnd = &t
	}
	meetings, err := c.conferenceService.AuthAndQueryMeetings(
		ctx,
		s.UserID,
		queryMeetingsRequest.Ids,
		queryMeetingsRequest.CreatorIds,
		queryMeetingsRequest.UserIds,
		queryMeetingsRequest.GroupIds,
		creationDateStart,
		creationDateEnd,
		queryMeetingsRequest.Skip,
		queryMeetingsRequest.Limit,
	)
	if err != nil {
		return nil, err
	}
	if len(meetings) == 0 {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(int32(constant.ResponseStatusCode_NO_CONTENT)),
		}, nil
	}
	protoMeetings := make([]*protocol.Meeting, len(meetings))
	for i, meeting := range meetings {
		protoMeetings[i] = &protocol.Meeting{
			Id:        meeting.ID,
			CreatorId: meeting.CreatorID,
			UserId:    meeting.UserID,
			GroupId:   meeting.GroupID,
			Name:      meeting.Name,
			Intro:     meeting.Intro,
			Password:  meeting.Password,
			StartDate: meeting.StartDate.UnixMilli(),
		}
		if !meeting.EndDate.IsZero() {
			protoMeetings[i].EndDate = proto.Int64(meeting.EndDate.UnixMilli())
		}
		if !meeting.CancelDate.IsZero() {
			protoMeetings[i].CancelDate = proto.Int64(meeting.CancelDate.UnixMilli())
		}
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Meetings{
				Meetings: &protocol.Meetings{
					Meetings: protoMeetings,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateMeetingRequest(@NotNull ClientRequest<UpdateMeetingRequest> clientRequest)
func (c *ConferenceServiceController) HandleUpdateMeetingRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateMeetingRequest := req.GetUpdateMeetingRequest()
	result, err := c.conferenceService.AuthAndUpdateMeeting(
		ctx,
		s.UserID,
		updateMeetingRequest.GetId(),
		updateMeetingRequest.Name,
		updateMeetingRequest.Intro,
		updateMeetingRequest.Password,
	)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, nil
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
	}, nil
}

// @MappedFrom handleUpdateMeetingInvitationRequest(@NotNull ClientRequest<UpdateMeetingInvitationRequest> clientRequest)
func (c *ConferenceServiceController) HandleUpdateMeetingInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateMeetingInvitationRequest := req.GetUpdateMeetingInvitationRequest()
	result, err := c.conferenceService.AuthAndUpdateMeetingInvitation(
		ctx,
		s.UserID,
		updateMeetingInvitationRequest.GetMeetingId(),
		updateMeetingInvitationRequest.Password,
		updateMeetingInvitationRequest.GetResponseAction(),
	)
	if err != nil {
		return nil, err
	}
	if !result.Updated {
		// Parity: Java version might return different codes depending on result.
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
	}, nil
}
