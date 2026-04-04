package dto

import "time"

// @MappedFrom AddFriendRequestDTO.java
type AddFriendRequestDTO struct {
	ID           *int64     `json:"id,omitempty"`
	RequesterID  *int64     `json:"requesterId,omitempty"`
	RecipientID  *int64     `json:"recipientId,omitempty"`
	Content      *string    `json:"content,omitempty"`
	Status       *int       `json:"status,omitempty"`
	Reason       *string    `json:"reason,omitempty"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
	ResponseDate *time.Time `json:"responseDate,omitempty"`
}

// @MappedFrom AddRelationshipDTO.java
type AddRelationshipDTO struct {
	OwnerID           *int64     `json:"ownerId,omitempty"`
	RelatedUserID     *int64     `json:"relatedUserId,omitempty"`
	Name              *string    `json:"name,omitempty"`
	BlockDate         *time.Time `json:"blockDate,omitempty"`
	EstablishmentDate *time.Time `json:"establishmentDate,omitempty"`
}

// @MappedFrom AddRelationshipGroupDTO.java
type AddRelationshipGroupDTO struct {
	OwnerID      *int64     `json:"ownerId,omitempty"`
	Index        *int       `json:"index,omitempty"`
	Name         *string    `json:"name,omitempty"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
}

// @MappedFrom AddUserDTO.java
type AddUserDTO struct {
	ID                *int64     `json:"id,omitempty"`
	Password          *string    `json:"password,omitempty"`
	Name              *string    `json:"name,omitempty"`
	Intro             *string    `json:"intro,omitempty"`
	ProfilePicture    *string    `json:"profilePicture,omitempty"`
	ProfileAccess     *int       `json:"profileAccess,omitempty"`
	PermissionGroupID *int64     `json:"permissionGroupId,omitempty"`
	RegistrationDate  *time.Time `json:"registrationDate,omitempty"`
	IsActive          *bool      `json:"isActive,omitempty"`
}

// @MappedFrom AddUserRoleDTO.java
type AddUserRoleDTO struct {
	ID                              *int64        `json:"id,omitempty"`
	Name                            *string       `json:"name,omitempty"`
	CreatableGroupTypeIDs           []int64       `json:"creatableGroupTypeIds,omitempty"`
	OwnedGroupLimit                 *int          `json:"ownedGroupLimit,omitempty"`
	OwnedGroupLimitForEachGroupType *int          `json:"ownedGroupLimitForEachGroupType,omitempty"`
	GroupTypeIDToLimit              map[int64]int `json:"groupTypeIdToLimit,omitempty"`
}

// @MappedFrom UpdateFriendRequestDTO.java
type UpdateFriendRequestDTO struct {
	RequesterID  *int64     `json:"requesterId,omitempty"`
	RecipientID  *int64     `json:"recipientId,omitempty"`
	Content      *string    `json:"content,omitempty"`
	Status       *int       `json:"status,omitempty"`
	Reason       *string    `json:"reason,omitempty"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
	ResponseDate *time.Time `json:"responseDate,omitempty"`
}

// @MappedFrom UpdateOnlineStatusDTO.java
type UpdateOnlineStatusDTO struct {
	OnlineStatus *int `json:"onlineStatus,omitempty"`
}

// @MappedFrom UpdateRelationshipDTO.java
type UpdateRelationshipDTO struct {
	Name              *string    `json:"name,omitempty"`
	BlockDate         *time.Time `json:"blockDate,omitempty"`
	EstablishmentDate *time.Time `json:"establishmentDate,omitempty"`
}

// @MappedFrom UpdateRelationshipGroupDTO.java
type UpdateRelationshipGroupDTO struct {
	Name         *string    `json:"name,omitempty"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
}

// @MappedFrom UpdateUserDTO.java
type UpdateUserDTO struct {
	Password          *string    `json:"password,omitempty"`
	Name              *string    `json:"name,omitempty"`
	Intro             *string    `json:"intro,omitempty"`
	ProfilePicture    *string    `json:"profilePicture,omitempty"`
	ProfileAccess     *int       `json:"profileAccess,omitempty"`
	PermissionGroupID *int64     `json:"permissionGroupId,omitempty"`
	RegistrationDate  *time.Time `json:"registrationDate,omitempty"`
	IsActive          *bool      `json:"isActive,omitempty"`
}

// @MappedFrom UpdateUserRoleDTO.java
type UpdateUserRoleDTO struct {
	Name                            *string       `json:"name,omitempty"`
	CreatableGroupTypeIDs           []int64       `json:"creatableGroupTypeIds,omitempty"`
	OwnedGroupLimit                 *int          `json:"ownedGroupLimit,omitempty"`
	OwnedGroupLimitForEachGroupType *int          `json:"ownedGroupLimitForEachGroupType,omitempty"`
	GroupTypeIDToLimit              map[int64]int `json:"groupTypeIdToLimit,omitempty"`
}

// @MappedFrom OnlineUserCountDTO.java
type OnlineUserCountDTO struct {
	Total             *int           `json:"total,omitempty"`
	NodeIDToUserCount map[string]int `json:"nodeIdToUserCount,omitempty"`
}

// @MappedFrom UserFriendRequestDTO.java
type UserFriendRequestDTO struct {
	ID             *int64     `json:"id,omitempty"`
	Content        *string    `json:"content,omitempty"`
	Status         *int       `json:"status,omitempty"`
	Reason         *string    `json:"reason,omitempty"`
	CreationDate   *time.Time `json:"creationDate,omitempty"`
	ResponseDate   *time.Time `json:"responseDate,omitempty"`
	RequesterID    *int64     `json:"requesterId,omitempty"`
	RecipientID    *int64     `json:"recipientId,omitempty"`
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
}

// @MappedFrom UserLocationDTO.java
type UserLocationDTO struct {
	UserID     *int64   `json:"userId,omitempty"`
	DeviceType *int     `json:"deviceType,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
}

// @MappedFrom UserRelationshipDTO.java
type UserRelationshipDTO struct {
	OwnerID           *int64     `json:"ownerId,omitempty"`
	RelatedUserID     *int64     `json:"relatedUserId,omitempty"`
	Name              *string    `json:"name,omitempty"`
	BlockDate         *time.Time `json:"blockDate,omitempty"`
	EstablishmentDate *time.Time `json:"establishmentDate,omitempty"`
	GroupIndexes      []int      `json:"groupIndexes,omitempty"`
}

// @MappedFrom UserStatisticsDTO.java
type UserStatisticsDTO struct {
	DeletedUsers             *int64        `json:"deletedUsers,omitempty"`
	UsersWhoSentMessages     *int64        `json:"usersWhoSentMessages,omitempty"`
	LoggedInUsers            *int64        `json:"loggedInUsers,omitempty"`
	MaxOnlineUsers           *int64        `json:"maxOnlineUsers,omitempty"`
	RegisteredUsers          *int64        `json:"registeredUsers,omitempty"`
	DeletedUsersRecords      []interface{} `json:"deletedUsersRecords,omitempty"`         // placeholder
	UsersWhoSentMessagesRecs []interface{} `json:"usersWhoSentMessagesRecords,omitempty"` // placeholder
	LoggedInUsersRecords     []interface{} `json:"loggedInUsersRecords,omitempty"`        // placeholder
	MaxOnlineUsersRecords    []interface{} `json:"maxOnlineUsersRecords,omitempty"`       // placeholder
	RegisteredUsersRecords   []interface{} `json:"registeredUsersRecords,omitempty"`      // placeholder
}

// @MappedFrom UserStatusDTO.java
type UserStatusDTO struct {
	UserID             *int64           `json:"userId,omitempty"`
	Status             *int             `json:"status,omitempty"`
	DeviceTypeToNodeID map[int]string   `json:"deviceTypeToNodeId,omitempty"`
	LoginDate          *time.Time       `json:"loginDate,omitempty"`
	LoginLocation      *UserLocationDTO `json:"loginLocation,omitempty"`
}
