package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
)

const GroupMemberCollectionName = "groupMember"

// @MappedFrom GroupMemberRepository
type GroupMemberRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewGroupMemberRepository(client *turmsmongo.Client) *GroupMemberRepository {
	return &GroupMemberRepository{
		client: client,
		col:    client.Collection(GroupMemberCollectionName),
	}
}

// AddGroupMember adds a member to a group or updates their role.
// @MappedFrom addGroupMember(@RequestBody AddGroupMemberDTO addGroupMemberDTO)
// @MappedFrom addGroupMember(@NotNull Long groupId, @NotNull Long userId, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)
func (r *GroupMemberRepository) AddGroupMember(ctx context.Context, member *po.GroupMember) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": member.ID}
	update := bson.M{"$set": member}

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// RemoveGroupMember removes a member from a group.
func (r *GroupMemberRepository) RemoveGroupMember(ctx context.Context, groupID, userID int64) error {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
	}
	_, err := r.col.DeleteOne(ctx, filter)
	return err
}

// FindGroupMemberRole retrieves the role of a user in a group.
// @MappedFrom findGroupMemberRole(Long userId, Long groupId)
func (r *GroupMemberRepository) FindGroupMemberRole(ctx context.Context, groupID, userID int64) (*protocol.GroupMemberRole, error) {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
	}
	opts := options.FindOne().SetProjection(bson.M{"role": 1})

	var member po.GroupMember
	if err := r.col.FindOne(ctx, filter, opts).Decode(&member); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User is not a member of the group
		}
		return nil, err
	}
	return &member.Role, nil
}

// FindGroupMemberIDs retrieves all user IDs within a group.
// @MappedFrom findGroupMemberIds(Long groupId)
// @MappedFrom findGroupMemberIds(Set<Long> groupIds)
func (r *GroupMemberRepository) FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	filter := bson.M{
		"_id.gid": groupID,
	}
	opts := options.Find().SetProjection(bson.M{"_id.uid": 1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}

	var userIDs []int64
	for _, m := range members {
		userIDs = append(userIDs, m.ID.UserID)
	}
	return userIDs, nil
}

// IsMemberMuted checks if a specific group member is currently muted.
// @MappedFrom isMemberMuted(@NotNull Long groupId, @NotNull Long userId, boolean preferCache)
// @MappedFrom isMemberMuted(Long groupId, Long userId)
func (r *GroupMemberRepository) IsMemberMuted(ctx context.Context, groupID, userID int64) (bool, error) {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
		"med": bson.M{"$gt": time.Now()},
	}

	count, err := r.col.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindUserJoinedGroupIDs retrieves all group IDs that a user belongs to.
// @MappedFrom findUserJoinedGroupIds(Long userId)
func (r *GroupMemberRepository) FindUserJoinedGroupIDs(ctx context.Context, userID int64) ([]int64, error) {
	filter := bson.M{
		"_id.uid": userID,
	}
	opts := options.Find().SetProjection(bson.M{"_id.gid": 1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}

	var groupIDs []int64
	for _, m := range members {
		groupIDs = append(groupIDs, m.ID.GroupID)
	}
	return groupIDs, nil
}

// DeleteByIds removes multiple group members by their keys.
func (r *GroupMemberRepository) DeleteByIds(ctx context.Context, keys []po.GroupMemberKey) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"_id": bson.M{"$in": keys},
	}
	return r.col.DeleteMany(ctx, filter)
}

// DeleteByGroupIDs removes all members of the specified groups.
func (r *GroupMemberRepository) DeleteByGroupIDs(ctx context.Context, groupIDs []int64) (*mongo.DeleteResult, error) {
	if len(groupIDs) == 0 {
		return &mongo.DeleteResult{}, nil
	}
	filter := bson.M{
		"_id.gid": bson.M{"$in": groupIDs},
	}
	return r.col.DeleteMany(ctx, filter)
}

// UpdateGroupMembers updates multiple group members' properties.
// @MappedFrom updateGroupMembers(List<GroupMember.Key> keys, @RequestBody UpdateGroupMemberDTO updateGroupMemberDTO)
// @MappedFrom updateGroupMembers(Set<GroupMember.Key> keys, @Nullable String name, @Nullable GroupMemberRole role, @Nullable Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)
// @MappedFrom updateGroupMembers(@NotNull Long groupId, @NotEmpty Set<Long> memberIds, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)
// @MappedFrom updateGroupMembers(@NotEmpty Set<GroupMember.Key> keys, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session, boolean updateGroupMembersVersion)
func (r *GroupMemberRepository) UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, name *string, role *protocol.GroupMemberRole, joinDate *time.Time, muteEndDate *time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"_id": bson.M{"$in": keys},
	}

	update := bson.M{}
	set := bson.M{}
	unset := bson.M{}

	if name != nil {
		set["n"] = *name
	}
	if role != nil {
		set["role"] = *role
	}
	if joinDate != nil {
		set["jd"] = *joinDate
	}
	if muteEndDate != nil {
		if muteEndDate.Before(time.Now()) {
			unset["med"] = ""
		} else {
			set["med"] = *muteEndDate
		}
	}

	if len(set) > 0 {
		update["$set"] = set
	}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	if len(update) == 0 {
		return &mongo.UpdateResult{MatchedCount: int64(len(keys))}, nil
	}

	return r.col.UpdateMany(ctx, filter, update)
}

// CountMembers returns the total number of members in a group.
// @MappedFrom countMembers(Long groupId)
// @MappedFrom countMembers(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes)
// @MappedFrom countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<@ValidGroupMemberRole GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)
// @MappedFrom countMembers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable Set<GroupMemberRole> roles, @Nullable DateRange joinDateRange, @Nullable DateRange muteEndDateRange)
func (r *GroupMemberRepository) CountMembers(ctx context.Context, groupID int64) (int64, error) {
	filter := bson.M{
		"_id.gid": groupID,
	}
	return r.col.CountDocuments(ctx, filter)
}

// FindGroupMemberKeyAndRolePairs retrieves the roles of multiple users in a group.
func (r *GroupMemberRepository) FindGroupMemberKeyAndRolePairs(ctx context.Context, groupID int64, userIDs []int64) ([]po.GroupMember, error) {
	filter := bson.M{
		"_id.gid": groupID,
		"_id.uid": bson.M{"$in": userIDs},
	}
	opts := options.Find().SetProjection(bson.M{"role": 1, "_id.uid": 1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// IsGroupMember checks if a user is a member of a group.
// @MappedFrom isGroupMember(@NotNull Long groupId, @NotNull Long userId, boolean preferCache)
// @MappedFrom isGroupMember(@NotEmpty Set<Long> groupIds, @NotNull Long userId)
func (r *GroupMemberRepository) IsGroupMember(ctx context.Context, groupID, userID int64) (bool, error) {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
	}
	count, err := r.col.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GroupMemberRepository) FindGroupManagersAndOwnerId(ctx context.Context, groupId int64) ([]po.GroupMember, error) {
	return nil, nil
}

// FindGroupMembers retrieves all members of a group.
// @MappedFrom findGroupMembers(Long groupId)
func (r *GroupMemberRepository) FindGroupMembers(ctx context.Context, groupID int64) ([]po.GroupMember, error) {
	filter := bson.M{
		"_id.gid": groupID,
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// FindGroupMembersWithIds retrieves specific members of a group.
// @MappedFrom findGroupMembers(Long groupId, Set<Long> memberIds)
func (r *GroupMemberRepository) FindGroupMembersWithIds(ctx context.Context, groupID int64, memberIDs []int64) ([]po.GroupMember, error) {
	filter := bson.M{
		"_id.gid": groupID,
		"_id.uid": bson.M{"$in": memberIDs},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

func (r *GroupMemberRepository) FindGroupsMembers(ctx context.Context, groupIds, userIds []int64, roles []int, joinDateRange, muteEndDateRange any, page, size *int) ([]po.GroupMember, error) {
	return nil, nil
}
func (r *GroupMemberRepository) FindGroupMemberKeyAndRoleParis(ctx context.Context, userIds []int64, groupId int64) ([]po.GroupMember, error) {
	return nil, nil
}
func (r *GroupMemberRepository) FindMemberIdsByGroupId(ctx context.Context, groupId int64) ([]int64, error) {
	return nil, nil
}
func (r *GroupMemberRepository) FindUsersJoinedGroupIds(ctx context.Context, groupIds, userIds []int64, page, size *int) ([]int64, error) {
	return nil, nil
}
