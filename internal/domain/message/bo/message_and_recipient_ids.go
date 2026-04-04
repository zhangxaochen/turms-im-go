package bo

import (
	"im.turms/server/internal/domain/message/po"
)

type MessageAndRecipientIDs struct {
	Message      *po.Message
	RecipientIDs []int64
}
