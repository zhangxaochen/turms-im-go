package controller

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/session"
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

	var resourceIdStr string
	if deleteReq.IdStr != nil {
		resourceIdStr = *deleteReq.IdStr
	} else if deleteReq.IdNum != nil {
		resourceIdStr = strconv.FormatInt(*deleteReq.IdNum, 10)
	}

	resType := mapStorageResourceType(deleteReq.Type)

	err := c.storageService.DeleteResource(ctx, s.UserID, resType, resourceIdStr)
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

	resType := mapStorageResourceType(uploadReq.Type)
	var name, mediaType string
	if uploadReq.Name != nil {
		name = *uploadReq.Name
	}
	if uploadReq.MediaType != nil {
		mediaType = *uploadReq.MediaType
	}

	// Assuming there's some custom logic or no max size provided in request, passing 0
	url, err := c.storageService.QueryResourceUploadInfo(ctx, s.UserID, resType, name, mediaType, 0)
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

// @MappedFrom handleQueryResourceDownloadInfoRequest()
func (c *StorageController) HandleQueryResourceDownloadInfoRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	downloadReq := req.GetQueryResourceDownloadInfoRequest()

	resType := mapStorageResourceType(downloadReq.Type)
	var resourceIdStr string
	if downloadReq.IdStr != nil {
		resourceIdStr = *downloadReq.IdStr
	} else if downloadReq.IdNum != nil {
		resourceIdStr = strconv.FormatInt(*downloadReq.IdNum, 10)
	}

	url, err := c.storageService.QueryResourceDownloadInfo(ctx, s.UserID, resType, resourceIdStr)
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
	// Features related to attachment sharing are not yet fully implemented in core
	return nil, errors.New("NotImplemented")
}

// @MappedFrom handleQueryMessageAttachmentInfosRequest()
func (c *StorageController) HandleQueryMessageAttachmentInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// Features related to attachment querying are not yet fully implemented in core
	return nil, errors.New("NotImplemented")
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
