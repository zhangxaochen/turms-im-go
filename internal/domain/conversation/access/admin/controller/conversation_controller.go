package controller

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/access/admin/controller"
	response_dto "im.turms/server/internal/domain/common/access/admin/dto/response"
	"im.turms/server/internal/domain/conversation/access/admin/dto"
	"im.turms/server/internal/domain/conversation/po"
	"im.turms/server/internal/domain/conversation/service"
	"im.turms/server/internal/infra/property"
)

// ConversationController maps to ConversationController.java
// @MappedFrom ConversationController
type ConversationController struct {
	*controller.BaseController
	conversationService *service.ConversationService
}

func NewConversationController(
	propertiesManager *property.TurmsPropertiesManager,
	conversationService *service.ConversationService,
) *ConversationController {
	return &ConversationController{
		BaseController:      controller.NewBaseController(propertiesManager),
		conversationService: conversationService,
	}
}

// @MappedFrom queryConversations(@QueryParam(required = false)
func (c *ConversationController) QueryConversations(
	ctx context.Context,
	privateConversationKeys []po.PrivateConversationKey,
	ownerIds []int64,
	groupIds []int64,
) (*dto.ConversationsDTO, error) {
	var privateConversations []*po.PrivateConversation
	var err error

	if len(privateConversationKeys) > 0 {
		privateConversations, err = c.conversationService.QueryPrivateConversations(ctx, privateConversationKeys)
		if err != nil {
			return nil, err
		}
	}
	if len(ownerIds) > 0 {
		morePrivateConversations, err := c.conversationService.QueryPrivateConversationsByOwnerIds(ctx, ownerIds)
		if err != nil {
			return nil, err
		}
		privateConversations = append(privateConversations, morePrivateConversations...)
	}

	var groupConversations []*po.GroupConversation
	if len(groupIds) > 0 {
		groupConversations, err = c.conversationService.QueryGroupConversations(ctx, groupIds)
		if err != nil {
			return nil, err
		}
	}

	return &dto.ConversationsDTO{
		PrivateConversations: privateConversations,
		GroupConversations:   groupConversations,
	}, nil
}

// @MappedFrom deleteConversations(@QueryParam(required = false)
func (c *ConversationController) DeleteConversations(
	ctx context.Context,
	privateConversationKeys []po.PrivateConversationKey,
	ownerIds []int64,
	groupIds []int64,
) (*response_dto.DeleteResultDTO, error) {
	totalDeleted := int64(0)
	if len(privateConversationKeys) > 0 {
		result, err := c.conversationService.DeletePrivateConversationsByKeys(ctx, privateConversationKeys)
		if err != nil {
			return nil, err
		}
		totalDeleted += result.DeletedCount
	}

	if len(ownerIds) > 0 {
		result, err := c.conversationService.DeletePrivateConversationsByUserIds(ctx, ownerIds)
		if err != nil {
			return nil, err
		}
		totalDeleted += result.DeletedCount
	}

	if len(groupIds) > 0 {
		result, err := c.conversationService.DeleteGroupConversations(ctx, groupIds)
		if err != nil {
			return nil, err
		}
		totalDeleted += result.DeletedCount
	}

	return &response_dto.DeleteResultDTO{DeletedCount: totalDeleted}, nil
}

// @MappedFrom updateConversations(@QueryParam(required = false)
func (c *ConversationController) UpdateConversations(
	ctx context.Context,
	privateConversationKeys []po.PrivateConversationKey,
	groupConversationMemberKeys []po.GroupConversionMemberKey,
	updateConversationDTO *dto.UpdateConversationDTO,
) error {
	if len(privateConversationKeys) > 0 {
		var readDate *time.Time
		if !updateConversationDTO.ReadDate.IsZero() {
			readDate = &updateConversationDTO.ReadDate
		}
		err := c.conversationService.UpsertPrivateConversationsReadDate(ctx, privateConversationKeys, readDate)
		if err != nil {
			return err
		}
	}

	if len(groupConversationMemberKeys) > 0 {
		var readDate *time.Time
		if !updateConversationDTO.ReadDate.IsZero() {
			readDate = &updateConversationDTO.ReadDate
		}
        
        // Convert to map for UpsertGroupConversationsReadDate
        groupIDToMemberIDs := make(map[int64][]int64)
        for _, key := range groupConversationMemberKeys {
            groupIDToMemberIDs[key.GroupID] = append(groupIDToMemberIDs[key.GroupID], key.MemberID)
        }

		err := c.conversationService.UpsertGroupConversationsReadDate(ctx, groupIDToMemberIDs, readDate)
		if err != nil {
			return err
		}
	}

	return nil
}
