package controller

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/common/dto"
	"im.turms/server/internal/domain/conversation/po"
	"im.turms/server/internal/domain/conversation/service"
	group_service "im.turms/server/internal/domain/group/service"
	sr_dto "im.turms/server/internal/domain/common/access/servicerequest/dto"
	"im.turms/server/internal/infra/property"
	"im.turms/server/pkg/protocol"
)

type ConversationServiceController struct {
	conversationService *service.ConversationService
	groupMemberService  *group_service.GroupMemberService
	propertiesManager   *property.TurmsPropertiesManager

	notifyRequesterOtherOnlineSessionsOfPrivateConversationReadDateUpdated bool
	notifyContactOfPrivateConversationReadDateUpdated                      bool

	notifyRequesterOtherOnlineSessionsOfGroupConversationReadDateUpdated bool
	notifyOtherGroupMembersOfGroupConversationReadDateUpdated            bool
}

func NewConversationServiceController(
	propertiesManager *property.TurmsPropertiesManager,
	conversationService *service.ConversationService,
	groupMemberService *group_service.GroupMemberService,
) *ConversationServiceController {
	c := &ConversationServiceController{
		propertiesManager:   propertiesManager,
		conversationService: conversationService,
		groupMemberService:  groupMemberService,
	}
	propertiesManager.NotifyAndAddGlobalPropertiesChangeListener(func(p *property.TurmsProperties) {
		n := p.Service.Notification
		c.notifyRequesterOtherOnlineSessionsOfPrivateConversationReadDateUpdated = n.PrivateConversationReadDateUpdated.NotifyRequesterOtherOnlineSessions
		c.notifyContactOfPrivateConversationReadDateUpdated = n.PrivateConversationReadDateUpdated.NotifyContact
		c.notifyRequesterOtherOnlineSessionsOfGroupConversationReadDateUpdated = n.GroupConversationReadDateUpdated.NotifyRequesterOtherOnlineSessions
		c.notifyOtherGroupMembersOfGroupConversationReadDateUpdated = n.GroupConversationReadDateUpdated.NotifyOtherGroupMembers
	})
	return c
}

// @MappedFrom handleQueryConversationsRequest()
func (c *ConversationServiceController) HandleQueryConversationsRequest(ctx context.Context, request *sr_dto.ServiceRequest) (*dto.RequestHandlerResult, error) {
	queryRequest := request.TurmsRequest.GetQueryConversationsRequest()
	targetIDs := queryRequest.GetTargetIds()
	groupIDs := queryRequest.GetGroupIds()
	// Bug fix: Return NO_CONTENT when both lists are empty (Java returns NO_CONTENT, not OK).
	if len(targetIDs) == 0 && len(groupIDs) == 0 {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_NO_CONTENT), nil
	}

	// Bug fix: Java queries either private OR group conversations, never both.
	// When targetIDs is non-empty, Java only queries private conversations (ignoring groupIDs).
	// When targetIDs is empty, Java falls through to query group conversations.
	if len(targetIDs) > 0 {
		keys := make([]po.PrivateConversationKey, len(targetIDs))
		for i, targetID := range targetIDs {
			// Bug fix: Swap OwnerID and TargetID to match Java's key construction.
			// Java: queryPrivateConversations(targetIds, clientRequest.userId())
			// creates Key(ownerId=requestTargetId, targetId=currentUserId).
			// Go was incorrectly: Key(OwnerID=currentUserId, TargetID=requestTargetId).
			keys[i] = po.PrivateConversationKey{
				OwnerID:  targetID,
				TargetID: request.UserId,
			}
		}
		privateConversations, err := c.conversationService.QueryPrivateConversations(ctx, keys)
		if err != nil {
			return nil, err
		}
		if len(privateConversations) == 0 {
			return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_NO_CONTENT), nil
		}
		return dto.RequestHandlerResultOfResponse(&protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Conversations{
				Conversations: c.conversations2proto(privateConversations, nil),
			},
		}), nil
	}

	// Group conversations
	groupConversations, err := c.conversationService.QueryGroupConversations(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	if len(groupConversations) == 0 {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_NO_CONTENT), nil
	}

	return dto.RequestHandlerResultOfResponse(&protocol.TurmsNotification_Data{
		Kind: &protocol.TurmsNotification_Data_Conversations{
			Conversations: c.conversations2proto(nil, groupConversations),
		},
	}), nil
}

// @MappedFrom handleUpdateTypingStatusRequest()
func (c *ConversationServiceController) HandleUpdateTypingStatusRequest(ctx context.Context, request *sr_dto.ServiceRequest) (*dto.RequestHandlerResult, error) {
	updateRequest := request.TurmsRequest.GetUpdateTypingStatusRequest()
	toId := updateRequest.ToId
	isGroupMessage := updateRequest.IsGroupMessage

	memberIds, err := c.conversationService.AuthAndUpdateTypingStatus(ctx, request.UserId, isGroupMessage, toId)
	if err != nil {
		return nil, err
	}
	return dto.RequestHandlerResultOfForwardRecipientsNotification(false, memberIds, request.TurmsRequest), nil
}

// @MappedFrom handleUpdateConversationRequest()
func (c *ConversationServiceController) HandleUpdateConversationRequest(ctx context.Context, request *sr_dto.ServiceRequest) (*dto.RequestHandlerResult, error) {
	updateRequest := request.TurmsRequest.GetUpdateConversationRequest()
	targetID := updateRequest.TargetId
	groupID := updateRequest.GroupId

	if targetID == nil && groupID == nil {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), nil
	}
	if targetID != nil && groupID != nil {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), nil
	}

	readDate := time.UnixMilli(updateRequest.GetReadDate())
	if targetID != nil {
		err := c.conversationService.AuthAndUpsertPrivateConversationReadDate(ctx, request.UserId, *targetID, &readDate)
		if err != nil {
			return nil, err
		}
		if c.notifyRequesterOtherOnlineSessionsOfPrivateConversationReadDateUpdated || c.notifyContactOfPrivateConversationReadDateUpdated {
			if c.notifyRequesterOtherOnlineSessionsOfPrivateConversationReadDateUpdated && c.notifyContactOfPrivateConversationReadDateUpdated {
				return dto.RequestHandlerResultOfForwardRecipientsNotification(true, []int64{*targetID}, request.TurmsRequest), nil
			} else if c.notifyRequesterOtherOnlineSessionsOfPrivateConversationReadDateUpdated {
				return dto.RequestHandlerResultOfForwardNotification(true, request.TurmsRequest), nil
			} else {
				return dto.RequestHandlerResultOfRecipientNotification(*targetID, request.TurmsRequest), nil
			}
		}
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_OK), nil
	}

	// Group conversation
	err := c.conversationService.AuthAndUpsertGroupConversationReadDate(ctx, *groupID, request.UserId, &readDate)
	if err != nil {
		return nil, err
	}

	if c.notifyRequesterOtherOnlineSessionsOfGroupConversationReadDateUpdated || c.notifyOtherGroupMembersOfGroupConversationReadDateUpdated {
		if c.notifyRequesterOtherOnlineSessionsOfGroupConversationReadDateUpdated && c.notifyOtherGroupMembersOfGroupConversationReadDateUpdated {
			memberIDs, err := c.groupMemberService.FindGroupMemberIDs(ctx, *groupID)
			if err != nil {
				return nil, err
			}
			return dto.RequestHandlerResultOfForwardRecipientsNotification(true, memberIDs, request.TurmsRequest), nil
		} else if c.notifyRequesterOtherOnlineSessionsOfGroupConversationReadDateUpdated {
			return dto.RequestHandlerResultOfForwardNotification(true, request.TurmsRequest), nil
		} else {
			// Bug fix: Java passes all member IDs (including requester) as recipients.
			// Go was incorrectly filtering out the requester from member IDs.
			memberIDs, err := c.groupMemberService.FindGroupMemberIDs(ctx, *groupID)
			if err != nil {
				return nil, err
			}
			return dto.RequestHandlerResultOfRecipientsNotification(memberIDs, request.TurmsRequest), nil
		}
	}

	return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_OK), nil
}

func (c *ConversationServiceController) conversations2proto(privateConversations []*po.PrivateConversation, groupConversations []*po.GroupConversation) *protocol.Conversations {
	protoPrivateConversations := make([]*protocol.PrivateConversation, len(privateConversations))
	for i, pc := range privateConversations {
		protoPrivateConversations[i] = &protocol.PrivateConversation{
			OwnerId:  pc.ID.OwnerID,
			TargetId: pc.ID.TargetID,
			ReadDate: pc.ReadDate.UnixMilli(),
		}
	}
	protoGroupConversations := make([]*protocol.GroupConversation, len(groupConversations))
	for i, gc := range groupConversations {
		memberReadDates := make(map[int64]int64, len(gc.MemberIDToReadDate))
		for k, v := range gc.MemberIDToReadDate {
			memberID, _ := strconv.ParseInt(k, 10, 64)
			memberReadDates[memberID] = v.UnixMilli()
		}
		protoGroupConversations[i] = &protocol.GroupConversation{
			GroupId:            gc.ID,
			MemberIdToReadDate: memberReadDates,
		}
	}
	return &protocol.Conversations{
		PrivateConversations: protoPrivateConversations,
		GroupConversations:   protoGroupConversations,
	}
}

// ConversationSettingsServiceController maps to ConversationSettingsServiceController.java
// @MappedFrom ConversationSettingsServiceController
type ConversationSettingsServiceController struct {
	service           *service.ConversationSettingsService
	propertiesManager *property.TurmsPropertiesManager
}

func NewConversationSettingsServiceController(
	service *service.ConversationSettingsService,
	propertiesManager *property.TurmsPropertiesManager,
) *ConversationSettingsServiceController {
	return &ConversationSettingsServiceController{
		service:           service,
		propertiesManager: propertiesManager,
	}
}

// @MappedFrom handleUpdateConversationSettingsRequest()
func (c *ConversationSettingsServiceController) HandleUpdateConversationSettingsRequest(ctx context.Context, request *sr_dto.ServiceRequest) (*dto.RequestHandlerResult, error) {
	updateRequest := request.TurmsRequest.GetUpdateConversationSettingsRequest()
	protoSettings := updateRequest.GetSettings()
	if len(protoSettings) == 0 {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_OK), nil
	}
	settings := make(map[string]any, len(protoSettings))
	for k, v := range protoSettings {
		settings[k] = protoValueToAny(v)
	}

	if updateRequest.GroupId != nil {
		updated, err := c.service.UpsertGroupConversationSettings(ctx, request.UserId, updateRequest.GetGroupId(), settings)
		if err != nil {
			return nil, err
		}
		if updated && c.propertiesManager.GetLocalProperties().Service.Notification.GroupConversationSettingUpdated.NotifyRequesterOtherOnlineSessions {
			return dto.RequestHandlerResultOfForwardNotification(true, request.TurmsRequest), nil
		}
	} else if updateRequest.UserId != nil {
		updated, err := c.service.UpsertPrivateConversationSettings(ctx, request.UserId, updateRequest.GetUserId(), settings)
		if err != nil {
			return nil, err
		}
		if updated && c.propertiesManager.GetLocalProperties().Service.Notification.PrivateConversationSettingUpdated.NotifyRequesterOtherOnlineSessions {
			return dto.RequestHandlerResultOfForwardNotification(true, request.TurmsRequest), nil
		}
	} else {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), nil
	}
	return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_OK), nil
}

// @MappedFrom handleDeleteConversationSettingsRequest()
func (c *ConversationSettingsServiceController) HandleDeleteConversationSettingsRequest(ctx context.Context, request *sr_dto.ServiceRequest) (*dto.RequestHandlerResult, error) {
	deleteRequest := request.TurmsRequest.GetDeleteConversationSettingsRequest()
	userIds := deleteRequest.GetUserIds()
	groupIds := deleteRequest.GetGroupIds()
	names := deleteRequest.GetNames()

	deleted, err := c.service.UnsetSettings(ctx, request.UserId, userIds, groupIds, names, nil)
	if err != nil {
		return nil, err
	}
	if deleted {
		hasUserId := len(userIds) > 0
		hasGroupId := len(groupIds) > 0
		if (hasUserId && c.propertiesManager.GetLocalProperties().Service.Notification.PrivateConversationSettingDeleted.NotifyRequesterOtherOnlineSessions) ||
			(hasGroupId && c.propertiesManager.GetLocalProperties().Service.Notification.GroupConversationSettingDeleted.NotifyRequesterOtherOnlineSessions) {
			return dto.RequestHandlerResultOfForwardNotification(true, request.TurmsRequest), nil
		}
	}
	return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_OK), nil
}

// @MappedFrom handleQueryConversationSettingsRequest()
func (c *ConversationSettingsServiceController) HandleQueryConversationSettingsRequest(ctx context.Context, request *sr_dto.ServiceRequest) (*dto.RequestHandlerResult, error) {
	queryRequest := request.TurmsRequest.GetQueryConversationSettingsRequest()
	userIds := queryRequest.GetUserIds()
	groupIds := queryRequest.GetGroupIds()
	names := queryRequest.GetNames()

	var lastUpdatedDateStart *time.Time
	if queryRequest.LastUpdatedDateStart != nil {
		t := time.UnixMilli(queryRequest.GetLastUpdatedDateStart())
		lastUpdatedDateStart = &t
	}

	poSettingsList, err := c.service.QuerySettings(ctx, request.UserId, userIds, groupIds, names, lastUpdatedDateStart)
	if err != nil {
		return nil, err
	}
	if len(poSettingsList) == 0 {
		return dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_NO_CONTENT), nil
	}

	protoSettingsList := make([]*protocol.ConversationSettings, len(poSettingsList))
	for i, s := range poSettingsList {
		protoSettingsList[i] = c.poSettingsToProto(s)
	}

	return dto.RequestHandlerResultOfResponse(&protocol.TurmsNotification_Data{
		Kind: &protocol.TurmsNotification_Data_ConversationSettingsList{
			ConversationSettingsList: &protocol.ConversationSettingsList{
				ConversationSettingsList: protoSettingsList,
			},
		},
	}), nil
}

func (c *ConversationSettingsServiceController) poSettingsToProto(s po.ConversationSettings) *protocol.ConversationSettings {
	protoSettings := &protocol.ConversationSettings{
		Settings:        make(map[string]*protocol.Value, len(s.Settings)),
		LastUpdatedDate: proto.Int64(s.LastUpdatedDate.UnixMilli()),
	}
	targetId := s.ID.TargetId
	if targetId < 0 {
		protoSettings.GroupId = proto.Int64(-targetId)
	} else {
		protoSettings.UserId = proto.Int64(targetId)
	}
	if len(s.Settings) > 0 {
		for k, v := range s.Settings {
			protoSettings.Settings[k] = anyToProtoValue(v)
		}
	}
	return protoSettings
}

func protoValueToAny(v *protocol.Value) any {
	if v == nil {
		return nil
	}
	if v.Kind != nil {
		switch val := v.Kind.(type) {
		case *protocol.Value_Int32Value:
			return int64(val.Int32Value)
		case *protocol.Value_Int64Value:
			return val.Int64Value
		case *protocol.Value_FloatValue:
			return val.FloatValue
		case *protocol.Value_DoubleValue:
			return val.DoubleValue
		case *protocol.Value_BoolValue:
			return val.BoolValue
		case *protocol.Value_BytesValue:
			return val.BytesValue
		case *protocol.Value_StringValue:
			return val.StringValue
		}
	}
	if len(v.ListValue) > 0 {
		list := make([]any, len(v.ListValue))
		for i, ev := range v.ListValue {
			list[i] = protoValueToAny(ev)
		}
		return list
	}
	return nil
}

func anyToProtoValue(v any) *protocol.Value {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case int:
		return &protocol.Value{Kind: &protocol.Value_Int32Value{Int32Value: int32(val)}}
	case int32:
		return &protocol.Value{Kind: &protocol.Value_Int32Value{Int32Value: val}}
	case int64:
		return &protocol.Value{Kind: &protocol.Value_Int64Value{Int64Value: val}}
	case float32:
		return &protocol.Value{Kind: &protocol.Value_FloatValue{FloatValue: val}}
	case float64:
		return &protocol.Value{Kind: &protocol.Value_DoubleValue{DoubleValue: val}}
	case bool:
		return &protocol.Value{Kind: &protocol.Value_BoolValue{BoolValue: val}}
	case []byte:
		return &protocol.Value{Kind: &protocol.Value_BytesValue{BytesValue: val}}
	case string:
		return &protocol.Value{Kind: &protocol.Value_StringValue{StringValue: val}}
	case []any:
		list := make([]*protocol.Value, len(val))
		for i, ev := range val {
			list[i] = anyToProtoValue(ev)
		}
		return &protocol.Value{ListValue: list}
	}
	return nil
}
