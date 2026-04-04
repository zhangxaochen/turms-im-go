package bo

import (
	"im.turms/server/internal/domain/user/po"
)

// HandleFriendRequestResult represents the result of handling a friend request.
// @MappedFrom HandleFriendRequestResult.java
type HandleFriendRequestResult struct {
	FriendRequest                          *po.UserFriendRequest
	NewGroupIndexForFriendRequestRequester *int32
	NewGroupIndexForFriendRequestRecipient *int32
}
