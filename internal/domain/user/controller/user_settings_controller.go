package controller

import (
	"context"

	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/pkg/protocol"
)

type UserSettingsController struct {
	userSettingsService service.UserSettingsService
}

func NewUserSettingsController(userSettingsService service.UserSettingsService) *UserSettingsController {
	return &UserSettingsController{
		userSettingsService: userSettingsService,
	}
}

// NOTE: Since protocol.TurmsRequest doesn't have UserSettings requests yet,
// we'll leave this empty for now.
func (c *UserSettingsController) RegisterRoutes(r *router.Router) {
	// The following will be enabled once protocol supports UserSettings
	// r.RegisterController(&protocol.TurmsRequest_UpdateUserSettingsRequest{}, c.HandleUpdateUserSettingsRequest)
	// r.RegisterController(&protocol.TurmsRequest_QueryUserSettingsRequest{}, c.HandleQueryUserSettingsRequest)
}

// HandleUpdateUserSettingsRequest updates user settings.
// @MappedFrom handleUpdateUserSettingsRequest()
func (c *UserSettingsController) HandleUpdateUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implementation depends on protocol.UpdateUserSettingsRequest
	return buildSuccessNotification(req.RequestId), nil
}

// HandleQueryUserSettingsRequest queries user settings.
// @MappedFrom handleQueryUserSettingsRequest()
func (c *UserSettingsController) HandleQueryUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implementation depends on protocol.QueryUserSettingsRequest
	return buildSuccessNotification(req.RequestId), nil
}
