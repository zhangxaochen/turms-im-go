package dto

import (
	"time"

	common_dto "im.turms/server/internal/domain/common/access/admin/dto"
)

// @MappedFrom AddGroupBlockedUserDTO
type AddGroupBlockedUserDTO struct {
	GroupId     *int64     `json:"groupId"`
	UserId      *int64     `json:"userId"`
	BlockDate   *time.Time `json:"blockDate"`
	RequesterId *int64     `json:"requesterId"`
}

// @MappedFrom AddGroupDTO
type AddGroupDTO struct {
	TypeId       *int64     `json:"typeId"`
	CreatorId    *int64     `json:"creatorId"`
	OwnerId      *int64     `json:"ownerId"`
	Name         *string    `json:"name"`
	Intro        *string    `json:"intro"`
	Announcement *string    `json:"announcement"`
	MinimumScore *int       `json:"minimumScore"`
	CreationDate *time.Time `json:"creationDate"`
	DeletionDate *time.Time `json:"deletionDate"`
	MuteEndDate  *time.Time `json:"muteEndDate"`
	IsActive     *bool      `json:"isActive"`
}

// @MappedFrom AddGroupInvitationDTO
type AddGroupInvitationDTO struct {
	Id           *int64     `json:"id"`
	Content      *string    `json:"content"`
	Status       any        `json:"status"` // RequestStatus
	CreationDate *time.Time `json:"creationDate"`
	ResponseDate *time.Time `json:"responseDate"`
	GroupId      *int64     `json:"groupId"`
	InviterId    *int64     `json:"inviterId"`
	InviteeId    *int64     `json:"inviteeId"`
}

// @MappedFrom AddGroupJoinQuestionDTO
type AddGroupJoinQuestionDTO struct {
	GroupId  *int64   `json:"groupId"`
	Question *string  `json:"question"`
	Answers  []string `json:"answers"`
	Score    *int     `json:"score"`
}

// @MappedFrom AddGroupJoinRequestDTO
type AddGroupJoinRequestDTO struct {
	Id             *int64     `json:"id"`
	Content        *string    `json:"content"`
	Status         any        `json:"status"`
	CreationDate   *time.Time `json:"creationDate"`
	ResponseDate   *time.Time `json:"responseDate"`
	ResponseReason *string    `json:"responseReason"`
	GroupId        *int64     `json:"groupId"`
	RequesterId    *int64     `json:"requesterId"`
	ResponderId    *int64     `json:"responderId"`
}

// @MappedFrom AddGroupMemberDTO
type AddGroupMemberDTO struct {
	GroupId     *int64     `json:"groupId"`
	UserId      *int64     `json:"userId"`
	Name        *string    `json:"name"`
	Role        any        `json:"role"` // GroupMemberRole
	JoinDate    *time.Time `json:"joinDate"`
	MuteEndDate *time.Time `json:"muteEndDate"`
}

// @MappedFrom AddGroupTypeDTO
type AddGroupTypeDTO struct {
	Name                     *string `json:"name"`
	GroupSizeLimit           *int    `json:"groupSizeLimit"`
	InvitationStrategy       any     `json:"invitationStrategy"`
	JoinStrategy             any     `json:"joinStrategy"`
	GroupInfoUpdateStrategy  any     `json:"groupInfoUpdateStrategy"`
	MemberInfoUpdateStrategy any     `json:"memberInfoUpdateStrategy"`
	GuestSpeakable           *bool   `json:"guestSpeakable"`
	SelfInfoUpdatable        *bool   `json:"selfInfoUpdatable"`
	EnableReadReceipt        *bool   `json:"enableReadReceipt"`
	MessageEditable          *bool   `json:"messageEditable"`
}

// @MappedFrom UpdateGroupBlockedUserDTO
type UpdateGroupBlockedUserDTO struct {
	BlockDate   *time.Time `json:"blockDate"`
	RequesterId *int64     `json:"requesterId"`
}

// @MappedFrom UpdateGroupDTO
type UpdateGroupDTO struct {
	TypeId            *int64     `json:"typeId"`
	CreatorId         *int64     `json:"creatorId"`
	OwnerId           *int64     `json:"ownerId"`
	Name              *string    `json:"name"`
	Intro             *string    `json:"intro"`
	Announcement      *string    `json:"announcement"`
	MinimumScore      *int       `json:"minimumScore"`
	IsActive          *bool      `json:"isActive"`
	CreationDate      *time.Time `json:"creationDate"`
	DeletionDate      *time.Time `json:"deletionDate"`
	MuteEndDate       *time.Time `json:"muteEndDate"`
	SuccessorId       *int64     `json:"successorId"`
	QuitAfterTransfer *bool      `json:"quitAfterTransfer"`
}

// @MappedFrom UpdateGroupInvitationDTO
type UpdateGroupInvitationDTO struct {
	Content      *string    `json:"content"`
	Status       any        `json:"status"`
	CreationDate *time.Time `json:"creationDate"`
	ResponseDate *time.Time `json:"responseDate"`
	GroupId      *int64     `json:"groupId"`
	InviterId    *int64     `json:"inviterId"`
	InviteeId    *int64     `json:"inviteeId"`
}

// @MappedFrom UpdateGroupJoinQuestionDTO
type UpdateGroupJoinQuestionDTO struct {
	GroupId  *int64   `json:"groupId"`
	Question *string  `json:"question"`
	Answers  []string `json:"answers"`
	Score    *int     `json:"score"`
}

// @MappedFrom UpdateGroupJoinRequestDTO
type UpdateGroupJoinRequestDTO struct {
	Content      *string    `json:"content"`
	Status       any        `json:"status"`
	CreationDate *time.Time `json:"creationDate"`
	ResponseDate *time.Time `json:"responseDate"`
	GroupId      *int64     `json:"groupId"`
	RequesterId  *int64     `json:"requesterId"`
	ResponderId  *int64     `json:"responderId"`
}

// @MappedFrom UpdateGroupMemberDTO
type UpdateGroupMemberDTO struct {
	Name        *string    `json:"name"`
	Role        any        `json:"role"`
	JoinDate    *time.Time `json:"joinDate"`
	MuteEndDate *time.Time `json:"muteEndDate"`
}

// @MappedFrom UpdateGroupTypeDTO
type UpdateGroupTypeDTO struct {
	Name                     *string `json:"name"`
	GroupSizeLimit           *int    `json:"groupSizeLimit"`
	InvitationStrategy       any     `json:"invitationStrategy"`
	JoinStrategy             any     `json:"joinStrategy"`
	GroupInfoUpdateStrategy  any     `json:"groupInfoUpdateStrategy"`
	MemberInfoUpdateStrategy any     `json:"memberInfoUpdateStrategy"`
	GuestSpeakable           *bool   `json:"guestSpeakable"`
	SelfInfoUpdatable        *bool   `json:"selfInfoUpdatable"`
	EnableReadReceipt        *bool   `json:"enableReadReceipt"`
	MessageEditable          *bool   `json:"messageEditable"`
}

// @MappedFrom GroupStatisticsDTO
type GroupStatisticsDTO struct {
	DeletedGroups                 *int64 `json:"deletedGroups"`
	GroupsThatSentMessages        *int64 `json:"groupsThatSentMessages"`
	CreatedGroups                 *int64 `json:"createdGroups"`
	DeletedGroupsRecords          []dto.StatisticsRecordDTO `json:"deletedGroupsRecords"`
	GroupsThatSentMessagesRecords []dto.StatisticsRecordDTO `json:"groupsThatSentMessagesRecords"`
	CreatedGroupsRecords          []dto.StatisticsRecordDTO `json:"createdGroupsRecords"`
}
