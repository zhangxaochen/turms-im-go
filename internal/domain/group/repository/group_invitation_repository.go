package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type GroupInvitationRepository interface {
	Insert(ctx context.Context, inv *po.GroupInvitation) error
	HasPendingInvitation(ctx context.Context, groupID, inviteeID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindInvitationsByInviteeID(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error)
	FindInvitationsByGroupID(ctx context.Context, groupID int64) ([]po.GroupInvitation, error)
}

type groupInvitationRepository struct {
	coll *mongo.Collection
}

func NewGroupInvitationRepository(client *turmsmongo.Client) GroupInvitationRepository {
	return &groupInvitationRepository{
		coll: client.Collection(po.CollectionNameGroupInvitation),
	}
}

func (r *groupInvitationRepository) Insert(ctx context.Context, inv *po.GroupInvitation) error {
	_, err := r.coll.InsertOne(ctx, inv)
	return err
}

func (r *groupInvitationRepository) HasPendingInvitation(ctx context.Context, groupID, inviteeID int64) (bool, error) {
	filter := bson.M{
		"gid":  groupID,
		"ieid": inviteeID,
		"stat": po.RequestStatusPending,
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupInvitationRepository) UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error) {
	filter := bson.M{
		"_id":  invitationID,
		"stat": po.RequestStatusPending,
	}
	updateOps := bson.M{
		"stat": newStatus,
		"rd":   responseDate,
	}
	if reason != nil {
		updateOps["rsn"] = *reason
	}
	update := bson.M{
		"$set": updateOps,
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *groupInvitationRepository) FindInvitationsByInviteeID(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"ieid": inviteeID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}

func (r *groupInvitationRepository) FindInvitationsByGroupID(ctx context.Context, groupID int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"gid": groupID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}
