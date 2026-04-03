package dto

// @MappedFrom CheckGroupQuestionAnswerResult
type CheckGroupQuestionAnswerResult struct {
	Joined      bool    `json:"joined"`
	GroupId     *int64  `json:"groupId"`
	QuestionIds []int64 `json:"questionIds"`
	Score       *int    `json:"score"`
}

// @MappedFrom HandleHandleGroupInvitationResult
type HandleHandleGroupInvitationResult struct {
	GroupInvitation           any  `json:"groupInvitation"`
	RequesterAddedAsNewMember bool `json:"requesterAddedAsNewMember"`
}

// @MappedFrom HandleHandleGroupJoinRequestResult
type HandleHandleGroupJoinRequestResult struct {
	GroupJoinRequest          any  `json:"groupJoinRequest"`
	RequesterAddedAsNewMember bool `json:"requesterAddedAsNewMember"`
}

// @MappedFrom NewGroupQuestion
type NewGroupQuestion struct {
	Question *string  `json:"question"`
	Answers  []string `json:"answers"`
	Score    *int     `json:"score"`
}
