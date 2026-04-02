package po

const CollectionNameUserRole = "userRole"

// UserRole represents a user's role and associated permissions.
// MongoDB Collection: userRole
type UserRole struct {
	ID                              int64            `bson:"_id"` // id
	Name                            string           `bson:"n"` // name
	CreatableGroupTypeIDs           []int64          `bson:"cgtid"` // creatableGroupTypeIds
	OwnedGroupLimit                 int32            `bson:"ogl"` // ownedGroupLimit
	OwnedGroupLimitForEachGroupType int32            `bson:"oglegt"` // ownedGroupLimitForEachGroupType
	GroupTypeIDToLimit              map[int64]int32  `bson:"gtl"` // groupTypeIdToLimit
}
