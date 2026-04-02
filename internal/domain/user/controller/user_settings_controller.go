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

// RegisterRoutes wires all UserSettings handlers to the gateway router.
// NOTE: Since protocol.TurmsRequest doesn't have UserSettings requests yet,
// we'll leave this empty for now or use placeholders if needed.
func (c *UserSettingsController) RegisterRoutes(r *router.Router) {
	// r.RegisterController(&protocol.TurmsRequest_UpdateUserSettingsRequest{}, c.HandleUpdateUserSettingsRequest)
	// r.RegisterController(&protocol.TurmsRequest_QueryUserSettingsRequest{}, c.HandleQueryUserSettingsRequest)
}

// HandleUpdateUserSettingsRequest updates user settings.
func (c *UserSettingsController) HandleUpdateUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// In the future, we'll parse the actual settings from the request.
	// For now, this is a placeholder for the logic.
	// settings := req.GetUpdateUserSettingsRequest().GetSettings() 
	
	// Example logic assuming we have user ID from session
	// err := c.userSettingsService.UpsertSettings(ctx, s.UserID, settings)
	// if err != nil {
	//     return nil, err
	// }
	
	return buildSuccessNotification(req.RequestId), nil
}

// HandleQueryUserSettingsRequest queries user settings.
func (c *UserSettingsController) HandleQueryUserSettingsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// queryReq := req.GetQueryUserSettingsRequest()
	// settings, err := c.userSettingsService.QuerySetting(ctx, s.UserID, queryReq.GetNames())
	
	// Placeholder return
	return buildSuccessNotification(req.RequestId), nil
}

