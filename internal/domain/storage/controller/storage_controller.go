package controller

import (
	"context"
	"errors"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/storage/bo"
	"im.turms/server/internal/domain/storage/constants"
	"im.turms/server/internal/domain/storage/service"
	"im.turms/server/pkg/protocol"
)

type StorageController struct {
	storageService *service.StorageService
}

func NewStorageController(storageService *service.StorageService) *StorageController {
	return &StorageController{
		storageService: storageService,
	}
}

// @MappedFrom handleDeleteResourceRequest()
func (c *StorageController) HandleDeleteResourceRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteResourceRequest()

	requesterID := s.UserID
	if requesterID == 0 {
		return nil, errors.New("ILLEGAL_ARGUMENT: requesterId must not be null")
	}

	resType := mapStorageResourceType(deleteReq.Type)
	if resType == 0 {
		return nil, errors.New("ILLEGAL_ARGUMENT: unrecognized storage resource type")
	}

	// Build resource ID string: prefer IdStr, fallback to IdNum
	var resourceIdStr string
	if deleteReq.IdStr != nil {
		resourceIdStr = *deleteReq.IdStr
	} else if deleteReq.IdNum != nil {
		resourceIdStr = strconv.FormatInt(*deleteReq.IdNum, 10)
	}

	// Resource-type-based validation (matching Java dispatch logic)
	switch resType {
	case constants.StorageResourceTypeGroupProfilePicture:
		if deleteReq.IdNum == nil {
			return nil, errors.New("ILLEGAL_ARGUMENT: The group ID must not be null")
		}
	case constants.StorageResourceTypeMessageAttachment:
		// resourceIdNum and resourceIdStr can both be nullable for MESSAGE_ATTACHMENT
	}

	err := c.storageService.DeleteResource(ctx, requesterID, resType, resourceIdStr)
	if err != nil {
		return nil, err
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
	}, nil
}

// @MappedFrom handleQueryResourceUploadInfoRequest()
func (c *StorageController) HandleQueryResourceUploadInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	uploadReq := req.GetQueryResourceUploadInfoRequest()

	requesterID := s.UserID
	if requesterID == 0 {
		return nil, errors.New("ILLEGAL_ARGUMENT: requesterId must not be null")
	}

	resType := mapStorageResourceType(uploadReq.Type)
	if resType == 0 {
		return nil, errors.New("ILLEGAL_ARGUMENT: unrecognized storage resource type")
	}

	var name, mediaType string
	if uploadReq.Name != nil {
		name = *uploadReq.Name
	}
	if uploadReq.MediaType != nil {
		mediaType = *uploadReq.MediaType
	}

	// Resource-type-based validation
	switch resType {
	case constants.StorageResourceTypeGroupProfilePicture:
		if uploadReq.IdNum == nil {
			return nil, errors.New("ILLEGAL_ARGUMENT: The group ID must not be null")
		}
	case constants.StorageResourceTypeMessageAttachment:
		// Java routes based on resourceIdNum:
		// null -> queryMessageAttachmentUploadInfo
		// negative -> queryMessageAttachmentUploadInfoInGroupConversation (with -resourceIdNum)
		// positive/zero -> queryMessageAttachmentUploadInfoInPrivateConversation
		// The Go service currently uses a single method, so we pass through.
	}

	// Pass resourceIdNum as part of the resource key for routing
	var resourceKey string
	if uploadReq.IdNum != nil {
		resourceKey = strconv.FormatInt(*uploadReq.IdNum, 10)
	}

	url, err := c.storageService.QueryResourceUploadInfo(ctx, requesterID, resType, uploadReq.IdNum, name, mediaType, 0, resourceKey)
	if err != nil {
		return nil, err
	}

	// Java returns Map<String, String> serialized as alternating key-value pairs.
	// We return the URL in a Strings array for protocol compatibility.
	strings := []string{url}
	if resourceKey != "" {
		strings = []string{resourceKey, url}
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_StringsWithVersion{
				StringsWithVersion: &protocol.StringsWithVersion{
					Strings: strings,
				},
			},
		},
	}, nil
}

// @MappedFrom handleQueryResourceDownloadInfoRequest()
func (c *StorageController) HandleQueryResourceDownloadInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	downloadReq := req.GetQueryResourceDownloadInfoRequest()

	requesterID := s.UserID
	if requesterID == 0 {
		return nil, errors.New("ILLEGAL_ARGUMENT: requesterId must not be null")
	}

	resType := mapStorageResourceType(downloadReq.Type)
	if resType == 0 {
		return nil, errors.New("ILLEGAL_ARGUMENT: unrecognized storage resource type")
	}

	var resourceIdStr string
	if downloadReq.IdStr != nil {
		resourceIdStr = *downloadReq.IdStr
	} else if downloadReq.IdNum != nil {
		resourceIdStr = strconv.FormatInt(*downloadReq.IdNum, 10)
	}

	// Resource-type-based validation
	switch resType {
	case constants.StorageResourceTypeUserProfilePicture:
		if downloadReq.IdNum == nil {
			return nil, errors.New("ILLEGAL_ARGUMENT: The user ID must not be null")
		}
	case constants.StorageResourceTypeGroupProfilePicture:
		if downloadReq.IdNum == nil {
			return nil, errors.New("ILLEGAL_ARGUMENT: The group ID must not be null")
		}
	}

	url, err := c.storageService.QueryResourceDownloadInfo(ctx, requesterID, resType, downloadReq.IdNum, resourceIdStr)
	if err != nil {
		return nil, err
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_StringsWithVersion{
				StringsWithVersion: &protocol.StringsWithVersion{
					Strings: []string{url},
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateMessageAttachmentInfoRequest()
func (c *StorageController) HandleUpdateMessageAttachmentInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateMessageAttachmentInfoRequest()
	requesterID := s.UserID

	var err error
	if updateReq.UserIdToShareWith != nil {
		err = c.storageService.ShareMessageAttachmentWithUser(ctx, requesterID, updateReq.AttachmentIdNum, updateReq.AttachmentIdStr, *updateReq.UserIdToShareWith)
	} else if updateReq.GroupIdToShareWith != nil {
		err = c.storageService.ShareMessageAttachmentWithGroup(ctx, requesterID, updateReq.AttachmentIdNum, updateReq.AttachmentIdStr, *updateReq.GroupIdToShareWith)
	} else if updateReq.UserIdToUnshareWith != nil {
		err = c.storageService.UnshareMessageAttachmentWithUser(ctx, requesterID, updateReq.AttachmentIdNum, updateReq.AttachmentIdStr, *updateReq.UserIdToUnshareWith)
	} else if updateReq.GroupIdToUnshareWith != nil {
		err = c.storageService.UnshareMessageAttachmentWithGroup(ctx, requesterID, updateReq.AttachmentIdNum, updateReq.AttachmentIdStr, *updateReq.GroupIdToUnshareWith)
	}

	if err != nil {
		return nil, err
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
	}, nil
}

// @MappedFrom handleQueryMessageAttachmentInfosRequest()
func (c *StorageController) HandleQueryMessageAttachmentInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryMessageAttachmentInfosRequest()
	requesterID := s.UserID

	var creationDateStart, creationDateEnd *time.Time
	if queryReq.CreationDateStart != nil {
		t := time.UnixMilli(*queryReq.CreationDateStart)
		creationDateStart = &t
	}
	if queryReq.CreationDateEnd != nil {
		t := time.UnixMilli(*queryReq.CreationDateEnd)
		creationDateEnd = &t
	}

	var infos []bo.StorageResourceInfo
	var err error

	if queryReq.InPrivateConversation != nil {
		if *queryReq.InPrivateConversation {
			infos, err = c.storageService.QueryMessageAttachmentInfosInPrivateConversations(
				ctx, requesterID, nil, creationDateStart, creationDateEnd, queryReq.AreSharedByMe)
		} else {
			var userIDs []int64
			if len(queryReq.UserIds) > 0 {
				userIDs = queryReq.UserIds
			}
			infos, err = c.storageService.QueryMessageAttachmentInfosInGroupConversations(
				ctx, requesterID, nil, userIDs, creationDateStart, creationDateEnd)
		}
	} else {
		hasUserIDs := len(queryReq.UserIds) > 0
		var userIDs []int64
		if hasUserIDs {
			userIDs = queryReq.UserIds
		}

		if len(queryReq.GroupIds) > 0 {
			infos, err = c.storageService.QueryMessageAttachmentInfosInGroupConversations(
				ctx, requesterID, queryReq.GroupIds, userIDs, creationDateStart, creationDateEnd)
		} else if hasUserIDs {
			infos, err = c.storageService.QueryMessageAttachmentInfosInPrivateConversations(
				ctx, requesterID, userIDs, creationDateStart, creationDateEnd, queryReq.AreSharedByMe)
		} else {
			infos, err = c.storageService.QueryMessageAttachmentInfosUploadedByRequester(
				ctx, requesterID, creationDateStart, creationDateEnd)
		}
	}

	if err != nil {
		return nil, err
	}

	// Java's storageResourceInfo2proto sets: idNum, idStr, mediaType, uploaderId, creationDate
	// It does NOT set the Name field on the proto.
	protoInfos := make([]*protocol.StorageResourceInfo, 0, len(infos))
	for _, info := range infos {
		protoInfo := &protocol.StorageResourceInfo{
			IdNum:      info.IDNum,
			IdStr:      info.IDStr,
			UploaderId: info.UploaderID,
		}
		// Only set MediaType if non-empty (Java: .setMediaType(info.mediaType()))
		if info.MediaType != "" {
			protoInfo.MediaType = proto.String(info.MediaType)
		}
		if !info.CreationDate.IsZero() {
			protoInfo.CreationDate = info.CreationDate.UnixMilli()
		}
		protoInfos = append(protoInfos, protoInfo)
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_StorageResourceInfos{
				StorageResourceInfos: &protocol.StorageResourceInfos{
					Infos: protoInfos,
				},
			},
		},
	}, nil
}

func mapStorageResourceType(protoType protocol.StorageResourceType) constants.StorageResourceType {
	switch protoType {
	case protocol.StorageResourceType_USER_PROFILE_PICTURE:
		return constants.StorageResourceTypeUserProfilePicture
	case protocol.StorageResourceType_GROUP_PROFILE_PICTURE:
		return constants.StorageResourceTypeGroupProfilePicture
	case protocol.StorageResourceType_MESSAGE_ATTACHMENT:
		return constants.StorageResourceTypeMessageAttachment
	default:
		return 0
	}
}
