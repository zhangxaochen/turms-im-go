package dispatcher

import (
	"context"

	"im.turms/server/internal/domain/common/dto"
)

// ClientRequestHandler maps to ClientRequestHandler.java
// @MappedFrom ClientRequestHandler
type ClientRequestHandler func(ctx context.Context, clientRequest *dto.ClientRequest) (*dto.RequestHandlerResult, error)
