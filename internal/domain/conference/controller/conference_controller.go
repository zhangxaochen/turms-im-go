package controller

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	commonservice "im.turms/server/internal/domain/common/service"
	"im.turms/server/internal/domain/conference/service"
	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/infra/property"
	"im.turms/server/pkg/protocol"
)

// ConferenceServiceController maps to ConferenceServiceController.java
// @MappedFrom ConferenceServiceController
type ConferenceServiceController struct {
	conferenceService      *service.ConferenceService
	outboundMessageService commonservice.OutboundMessageService
	propertiesManager      *property.TurmsPropertiesManager
}

func NewConferenceServiceController(
	conferenceService *service.ConferenceService,
	outboundMessageService commonservice.OutboundMessageService,
	propertiesManager *property.TurmsPropertiesManager,
) *ConferenceServiceController {
	return &ConferenceServiceController{
		conferenceService:      conferenceService,
		outboundMessageService: outboundMessageService,
		propertiesManager:      propertiesManager,
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

	p := c.propertiesManager.GetLocalProperties()
	if p.Service.Notification.MeetingCanceled.NotifyMeetingParticipants {
		participantIds, err := c.conferenceService.QueryMeetingParticipants(ctx, result.Meeting.UserID, result.Meeting.GroupID)
		if err == nil && len(participantIds) > 0 {
			filteredIds := make([]int64, 0, len(participantIds))
			for _, id := range participantIds {
				if id != s.UserID {
					filteredIds = append(filteredIds, id)
				}
			}
			if len(filteredIds) > 0 && c.outboundMessageService != nil {
				notif := &protocol.TurmsNotification{
					RelayedRequest: req,
				}
				_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notif, filteredIds)
			}
		}
	} else if p.Service.Notification.MeetingCanceled.NotifyRequesterOtherOnlineSessions {
		if c.outboundMessageService != nil {
			notif := &protocol.TurmsNotification{
				RelayedRequest: req,
			}
			_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notif, []int64{s.UserID}) // Forward to other sessions
		}
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
	
	ids := queryMeetingsRequest.Ids
	creatorIds := queryMeetingsRequest.CreatorIds
	userIds := queryMeetingsRequest.UserIds
	groupIds := queryMeetingsRequest.GroupIds
	
	meetings, err := c.conferenceService.AuthAndQueryMeetings(
		ctx,
		s.UserID,
		ids,
		creatorIds,
		userIds,
		groupIds,
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

	p := c.propertiesManager.GetLocalProperties()
	if p.Service.Notification.MeetingUpdated.NotifyMeetingParticipants && (updateMeetingRequest.Name != nil || updateMeetingRequest.Intro != nil) {
		participantIds, err := c.conferenceService.QueryMeetingParticipants(ctx, result.Meeting.UserID, result.Meeting.GroupID)
		if err == nil && len(participantIds) > 0 {
			filteredIds := make([]int64, 0, len(participantIds))
			for _, id := range participantIds {
				if id != s.UserID {
					filteredIds = append(filteredIds, id)
				}
			}
			if len(filteredIds) > 0 && c.outboundMessageService != nil {
				// Clear password before sending
				copiedReq := proto.Clone(req).(*protocol.TurmsRequest)
				copiedReq.GetUpdateMeetingRequest().Password = nil
				notif := &protocol.TurmsNotification{
					RelayedRequest: copiedReq,
				}
				_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notif, filteredIds)
			}
		}
	} else if p.Service.Notification.MeetingUpdated.NotifyRequesterOtherOnlineSessions && (updateMeetingRequest.Name != nil || updateMeetingRequest.Intro != nil) {
		if c.outboundMessageService != nil {
			notif := &protocol.TurmsNotification{
				RelayedRequest: req,
			}
			_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notif, []int64{s.UserID})
		}
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

	var data *protocol.TurmsNotification_Data
	if result.Updated && updateMeetingInvitationRequest.GetResponseAction() == protocol.ResponseAction_ACCEPT && result.AccessToken != nil {
		data = &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_String_{
				String_: *result.AccessToken,
			},
		}
	}

	p := c.propertiesManager.GetLocalProperties()
	if result.Updated {
		if p.Service.Notification.MeetingInvitationUpdated.NotifyMeetingParticipants {
			participantIds, err := c.conferenceService.QueryMeetingParticipants(ctx, result.Meeting.UserID, result.Meeting.GroupID)
			if err == nil && len(participantIds) > 0 {
				filteredIds := make([]int64, 0, len(participantIds))
				for _, id := range participantIds {
					if id != s.UserID {
						filteredIds = append(filteredIds, id)
					}
				}
				if len(filteredIds) > 0 && c.outboundMessageService != nil {
					copiedReq := proto.Clone(req).(*protocol.TurmsRequest)
					copiedReq.GetUpdateMeetingInvitationRequest().Password = nil
					notif := &protocol.TurmsNotification{
						RelayedRequest: copiedReq,
					}
					_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notif, filteredIds)
				}
			}
		} else if p.Service.Notification.MeetingInvitationUpdated.NotifyRequesterOtherOnlineSessions {
			if c.outboundMessageService != nil {
				notif := &protocol.TurmsNotification{
					RelayedRequest: req,
				}
				_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notif, []int64{s.UserID})
			}
		}
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
		Data:      data,
	}, nil
}
