package dto

import (
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/pkg/protocol"
)

// Notification maps to nested Notification class in RequestHandlerResult
// @MappedFrom Notification
type Notification struct {
	ForwardToRequesterOtherOnlineSessions bool
	Recipients                            []int64
	Notification                          *protocol.TurmsRequest
}

// @MappedFrom Notification(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)
func NewNotification(forwardToRequesterOtherOnlineSessions bool, recipients []int64, notification *protocol.TurmsRequest) *Notification {
	return &Notification{
		ForwardToRequesterOtherOnlineSessions: forwardToRequesterOtherOnlineSessions,
		Recipients:                            recipients,
		Notification:                          notification,
	}
}

// RequestHandlerResult maps to RequestHandlerResult in Java.
// @MappedFrom RequestHandlerResult
type RequestHandlerResult struct {
	Code          constant.ResponseStatusCode
	Reason        *string
	Response      *protocol.TurmsNotification_Data
	Notifications []*Notification
}

// @MappedFrom RequestHandlerResult(ResponseStatusCode code, @Nullable String reason, @Nullable TurmsNotification.Data response, List<Notification> notifications)
func NewRequestHandlerResult(code constant.ResponseStatusCode, reason *string, response *protocol.TurmsNotification_Data, notifications []*Notification) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:          code,
		Reason:        reason,
		Response:      response,
		Notifications: notifications,
	}
}

// factory methods for RequestHandlerResult

// @MappedFrom of(@NotNull ResponseStatusCode code)
func RequestHandlerResultOfCode(code constant.ResponseStatusCode) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull ResponseStatusCode code, @Nullable String reason)
func RequestHandlerResultOfCodeReason(code constant.ResponseStatusCode, reason *string) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull TurmsNotification.Data response)
func RequestHandlerResultOfResponse(response *protocol.TurmsNotification_Data) *RequestHandlerResult { return nil }

// @MappedFrom of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)
func RequestHandlerResultOfForwardNotification(forward bool, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)
func RequestHandlerResultOfForwardRecipientNotification(forward bool, recipientId int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull Long recipientId, @NotNull TurmsRequest notification)
func RequestHandlerResultOfRecipientNotification(recipientId int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest dataForRecipient)
func RequestHandlerResultOfRecipientsNotification(recipientIds []int64, dataForRecipient *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)
func RequestHandlerResultOfForwardRecipientsNotification(forward bool, recipientIds []int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(TurmsNotification.Data response, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)
func RequestHandlerResultOfResponseRecipientsNotifications(response *protocol.TurmsNotification_Data, recipientIds []int64, notificationForRecipients, notificationForRequesterOtherOnlineSessions *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)
func RequestHandlerResultOfRecipientsNotifications(recipientIds []int64, notificationForRecipients, notificationForRequesterOtherOnlineSessions *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(TurmsNotification.Data response, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)
func RequestHandlerResultOfResponseForwardRecipientsNotification(response *protocol.TurmsNotification_Data, forward bool, recipientIds []int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull ResponseStatusCode code, @NotNull Long recipientId, @NotNull TurmsRequest notification)
func RequestHandlerResultOfCodeRecipientNotification(code constant.ResponseStatusCode, recipientId int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull ResponseStatusCode code, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)
func RequestHandlerResultOfCodeForwardRecipientNotification(code constant.ResponseStatusCode, forward bool, recipientId int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull List<Notification> notifications)
func RequestHandlerResultOfNotifications(notifications []*Notification) *RequestHandlerResult { return nil }

// @MappedFrom of(@NotNull Notification notification)
func RequestHandlerResultOfNotification(notification *Notification) *RequestHandlerResult { return nil }


// @MappedFrom ofDataLong(@NotNull Long value)
func RequestHandlerResultOfDataLong(value int64) *RequestHandlerResult { return nil }

// @MappedFrom ofDataLong(@NotNull Long value, @NotNull Long recipientId, @NotNull TurmsRequest notification)
func RequestHandlerResultOfDataLongRecipientNotification(value int64, recipientId int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom ofDataLong(@NotNull Long value, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)
func RequestHandlerResultOfDataLongForwardRecipientNotification(value int64, forward bool, recipientId int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom ofDataLong(@NotNull Long value, boolean forwardDataForRecipientsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)
func RequestHandlerResultOfDataLongForwardNotification(value int64, forward bool, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom ofDataLong(@NotNull Long value, boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipients, TurmsRequest notification)
func RequestHandlerResultOfDataLongForwardRecipientsNotification(value int64, forward bool, recipients []int64, notification *protocol.TurmsRequest) *RequestHandlerResult { return nil }

// @MappedFrom ofDataLongs(@NotNull Collection<Long> values)
func RequestHandlerResultOfDataLongs(values []int64) *RequestHandlerResult { return nil }


// factory methods for Notification

// @MappedFrom of(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)
func NotificationOfForwardRecipientsNotification(forward bool, recipients []int64, notification *protocol.TurmsRequest) *Notification { return nil }

// @MappedFrom of(boolean forwardToRequesterOtherOnlineSessions, Long recipient, TurmsRequest notification)
func NotificationOfForwardRecipientNotification(forward bool, recipient int64, notification *protocol.TurmsRequest) *Notification { return nil }

// @MappedFrom of(boolean forwardToRequesterOtherOnlineSessions, TurmsRequest notification)
func NotificationOfForwardNotification(forward bool, notification *protocol.TurmsRequest) *Notification { return nil }
