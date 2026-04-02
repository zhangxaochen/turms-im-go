package po

import "im.turms/server/internal/domain/group/constant"

// GroupType represents a group type entity in MongoDB (groupType collection).
// No need to shard because there are only a few/some group types.
type GroupType struct {
	ID                       int64                            `bson:"_id"`
	Name                     string                           `bson:"n"`
	GroupSizeLimit           int32                            `bson:"gsl"`
	InvitationStrategy       constant.GroupInvitationStrategy `bson:"is"`
	JoinStrategy             constant.GroupJoinStrategy       `bson:"js"`
	GroupInfoUpdateStrategy  constant.GroupUpdateStrategy     `bson:"gius"`
	MemberInfoUpdateStrategy constant.GroupUpdateStrategy     `bson:"mius"`
	GuestSpeakable           bool                             `bson:"gs"`
	SelfInfoUpdatable        bool                             `bson:"siu"`
	EnableReadReceipt        bool                             `bson:"err"`
	MessageEditable          bool                             `bson:"me"`
}
