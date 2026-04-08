package dto

import (
	"encoding/json"
	"fmt"
	"time"

	common_dto "im.turms/server/internal/domain/common/access/admin/dto"
)

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
// password is allowed for deserialization but excluded from serialization (matching Java @SensitiveProperty).
type AddUserDTO struct {
	ID                    *int64     `json:"id,omitempty"`
	password              *string    `json:"-"` // internal only; see UnmarshalJSON/MarshalJSON
	Name                  *string    `json:"name,omitempty"`
	Intro                 *string    `json:"intro,omitempty"`
	ProfilePicture        *string    `json:"profilePicture,omitempty"`
	ProfileAccessStrategy *int       `json:"profileAccessStrategy,omitempty"`
	RoleID                *int64     `json:"roleId,omitempty"`
	RegistrationDate      *time.Time `json:"registrationDate,omitempty"`
	IsActive              *bool      `json:"isActive,omitempty"`
}

// Password returns the password field (for internal use only, never serialized to JSON).
func (dto *AddUserDTO) Password() *string {
	return dto.password
}

// UnmarshalJSON implements custom deserialization to handle the password field
// which is allowed in input but excluded from output (matching Java @SensitiveProperty(ALLOW_DESERIALIZATION)).
func (dto *AddUserDTO) UnmarshalJSON(data []byte) error {
	type Alias AddUserDTO
	aux := &struct {
		Password *string `json:"password,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(dto),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	dto.password = aux.Password
	return nil
}

// MarshalJSON implements custom serialization to exclude the password field from JSON output.
func (dto *AddUserDTO) MarshalJSON() ([]byte, error) {
	type Alias AddUserDTO
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(dto),
	})
}

// String implements fmt.Stringer with password masking.
// @MappedFrom AddUserDTO.toString()
func (dto *AddUserDTO) String() string {
	return fmt.Sprintf(
		"AddUserDTO{id=%v, password=***, name=%v, intro=%v, profilePicture=%v, profileAccessStrategy=%v, roleId=%v, registrationDate=%v, isActive=%v}",
		dto.ID, dto.Name, dto.Intro, dto.ProfilePicture, dto.ProfileAccessStrategy, dto.RoleID, dto.RegistrationDate, dto.IsActive,
	)
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
// password is allowed for deserialization but excluded from serialization (matching Java @SensitiveProperty).
type UpdateUserDTO struct {
	password              *string    `json:"-"` // internal only; see UnmarshalJSON/MarshalJSON
	Name                  *string    `json:"name,omitempty"`
	Intro                 *string    `json:"intro,omitempty"`
	ProfilePicture        *string    `json:"profilePicture,omitempty"`
	ProfileAccessStrategy *int       `json:"profileAccessStrategy,omitempty"`
	RoleID                *int64     `json:"roleId,omitempty"`
	RegistrationDate      *time.Time `json:"registrationDate,omitempty"`
	IsActive              *bool      `json:"isActive,omitempty"`
}

// Password returns the password field (for internal use only, never serialized to JSON).
func (dto *UpdateUserDTO) Password() *string {
	return dto.password
}

// UnmarshalJSON implements custom deserialization to handle the password field.
func (dto *UpdateUserDTO) UnmarshalJSON(data []byte) error {
	type Alias UpdateUserDTO
	aux := &struct {
		Password *string `json:"password,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(dto),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	dto.password = aux.Password
	return nil
}

// MarshalJSON implements custom serialization to exclude the password field from JSON output.
func (dto *UpdateUserDTO) MarshalJSON() ([]byte, error) {
	type Alias UpdateUserDTO
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(dto),
	})
}

// String implements fmt.Stringer with password masking.
// @MappedFrom UpdateUserDTO.toString()
func (dto *UpdateUserDTO) String() string {
	return fmt.Sprintf(
		"UpdateUserDTO{password=***, name=%v, intro=%v, profilePicture=%v, profileAccessStrategy=%v, roleId=%v, registrationDate=%v, isActive=%v}",
		dto.Name, dto.Intro, dto.ProfilePicture, dto.ProfileAccessStrategy, dto.RoleID, dto.RegistrationDate, dto.IsActive,
	)
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

// UserRelationshipDTOKey represents the nested Key record from Java's UserRelationshipDTO.Key.
// @MappedFrom UserRelationshipDTO.Key
type UserRelationshipDTOKey struct {
	OwnerID       *int64 `json:"ownerId,omitempty"`
	RelatedUserID *int64 `json:"relatedUserId,omitempty"`
}

// @MappedFrom UserRelationshipDTO.java
type UserRelationshipDTO struct {
	Key               *UserRelationshipDTOKey `json:"key,omitempty"`
	Name              *string                 `json:"name,omitempty"`
	BlockDate         *time.Time              `json:"blockDate,omitempty"`
	EstablishmentDate *time.Time              `json:"establishmentDate,omitempty"`
	GroupIndexes      []int                   `json:"groupIndexes,omitempty"`
}

// @MappedFrom UserStatisticsDTO.java
type UserStatisticsDTO struct {
	DeletedUsers                 *int64                       `json:"deletedUsers,omitempty"`
	UsersWhoSentMessages         *int64                       `json:"usersWhoSentMessages,omitempty"`
	LoggedInUsers                *int64                       `json:"loggedInUsers,omitempty"`
	MaxOnlineUsers               *int64                       `json:"maxOnlineUsers,omitempty"`
	RegisteredUsers              *int64                       `json:"registeredUsers,omitempty"`
	DeletedUsersRecords          []common_dto.StatisticsRecordDTO `json:"deletedUsersRecords,omitempty"`
	UsersWhoSentMessagesRecords  []common_dto.StatisticsRecordDTO `json:"usersWhoSentMessagesRecords,omitempty"`
	LoggedInUsersRecords         []common_dto.StatisticsRecordDTO `json:"loggedInUsersRecords,omitempty"`
	MaxOnlineUsersRecords        []common_dto.StatisticsRecordDTO `json:"maxOnlineUsersRecords,omitempty"`
	RegisteredUsersRecords       []common_dto.StatisticsRecordDTO `json:"registeredUsersRecords,omitempty"`
}

// @MappedFrom UserStatusDTO.java
type UserStatusDTO struct {
	UserID             *int64             `json:"userId,omitempty"`
	Status             *int               `json:"status,omitempty"`
	DeviceTypeToNodeID map[int]string     `json:"deviceTypeToNodeId,omitempty"`
	LoginDate          *time.Time         `json:"loginDate,omitempty"`
	LoginLocation      *LocationDTO       `json:"loginLocation,omitempty"`
}

// LocationDTO represents the Java Location record with longitude, latitude, timestamp, and details.
// @MappedFrom Location.java (im.turms.server.common.domain.user.po.Location)
type LocationDTO struct {
	Longitude *float32          `json:"longitude,omitempty"`
	Latitude  *float32          `json:"latitude,omitempty"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}
