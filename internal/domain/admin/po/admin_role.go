package po

import (
	"time"

	"im.turms/server/internal/domain/admin/permission"
)

const CollectionNameAdminRole = "adminRole"

type AdminRole struct {
	ID           int64                        `bson:"_id"`
	Name         string                       `bson:"n"`
	Permissions  []permission.AdminPermission `bson:"perm"`
	Rank         int                          `bson:"rank"`
	CreationDate time.Time                    `bson:"cd"`
}

// Fields mapping for AdminRole
const (
	AdminRoleFieldID           = "_id"
	AdminRoleFieldName         = "n"
	AdminRoleFieldPermissions  = "perm"
	AdminRoleFieldRank         = "rank"
	AdminRoleFieldCreationDate = "cd"
)
