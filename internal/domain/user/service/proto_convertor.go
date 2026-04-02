package service

import (
	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/pkg/protocol"
)

func FriendRequestToProto(r *po.UserFriendRequest) *protocol.UserFriendRequest {
	if r == nil {
		return nil
	}
	var reason *string
	if r.Reason != nil {
		reason = r.Reason
	}
	status := func() protocol.RequestStatus {
		switch r.Status {
		case po.RequestStatusAccepted:
			return protocol.RequestStatus_ACCEPTED
		case po.RequestStatusDeclined:
			return protocol.RequestStatus_DECLINED
		case po.RequestStatusIgnored:
			return protocol.RequestStatus_IGNORED
		case po.RequestStatusExpired:
			return protocol.RequestStatus_EXPIRED
		case po.RequestStatusCanceled:
			return protocol.RequestStatus_CANCELED
		default:
			return protocol.RequestStatus_PENDING
		}
	}()
	return &protocol.UserFriendRequest{
		Id:            proto.Int64(r.ID),
		CreationDate:  proto.Int64(r.CreationDate.UnixMilli()),
		Content:       proto.String(r.Content),
		RequestStatus: &status,
		Reason:        reason,
		RequesterId:   proto.Int64(r.RequesterID),
		RecipientId:   proto.Int64(r.RecipientID),
	}
}

func RelationshipToProto(r *po.UserRelationship) *protocol.UserRelationship {
	if r == nil {
		return nil
	}
	var blockDate *int64
	if r.BlockDate != nil {
		t := r.BlockDate.UnixMilli()
		blockDate = &t
	}
	var establishmentDate *int64
	if r.EstablishmentDate != nil {
		t := r.EstablishmentDate.UnixMilli()
		establishmentDate = &t
	}
	var groupIndex *int64
	if r.GroupIndex != nil {
		idx := int64(*r.GroupIndex)
		groupIndex = &idx
	}
	return &protocol.UserRelationship{
		OwnerId:           proto.Int64(r.ID.OwnerID),
		RelatedUserId:     proto.Int64(r.ID.RelatedUserID),
		BlockDate:         blockDate,
		GroupIndex:        groupIndex,
		EstablishmentDate: establishmentDate,
	}
}

func RelationshipGroupToProto(g *po.UserRelationshipGroup) *protocol.UserRelationshipGroup {
	if g == nil {
		return nil
	}
	return &protocol.UserRelationshipGroup{
		Index: g.Key.Index,
		Name:  g.Name,
	}
}
