package constant

// GroupInvitationStrategy defines the strategy for sending group invitations.
type GroupInvitationStrategy int32

const (
	GroupInvitationStrategy_ALL_REQUIRING_APPROVAL                  GroupInvitationStrategy = 0
	GroupInvitationStrategy_OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL GroupInvitationStrategy = 1
	GroupInvitationStrategy_OWNER_MANAGER_REQUIRING_APPROVAL        GroupInvitationStrategy = 2
	GroupInvitationStrategy_OWNER_REQUIRING_APPROVAL                GroupInvitationStrategy = 3
	GroupInvitationStrategy_ALL                                     GroupInvitationStrategy = 4
	GroupInvitationStrategy_OWNER_MANAGER_MEMBER                    GroupInvitationStrategy = 5
	GroupInvitationStrategy_OWNER_MANAGER                           GroupInvitationStrategy = 6
	GroupInvitationStrategy_OWNER                                   GroupInvitationStrategy = 7
)

// RequiresApproval returns true if the strategy requires the invitee's approval.
func (s GroupInvitationStrategy) RequiresApproval() bool {
	return s <= GroupInvitationStrategy_OWNER_REQUIRING_APPROVAL
}

// GroupJoinStrategy defines the strategy for joining a group.
type GroupJoinStrategy int32

const (
	// Add the requester as a group member when the server received a membership request.
	GroupJoinStrategy_MEMBERSHIP_REQUEST GroupJoinStrategy = 0
	// A user can only join these groups via invitations.
	GroupJoinStrategy_INVITATION GroupJoinStrategy = 1
	// A user is required to answer questions to join.
	GroupJoinStrategy_QUESTION GroupJoinStrategy = 2
	// A user sends a join request to the server, and can only join the group automatically when the request is approved.
	GroupJoinStrategy_JOIN_REQUEST GroupJoinStrategy = 3
)

// GroupUpdateStrategy defines the strategy for updating group/member info.
type GroupUpdateStrategy int32

const (
	GroupUpdateStrategy_ALL                  GroupUpdateStrategy = 0
	GroupUpdateStrategy_OWNER_MANAGER_MEMBER GroupUpdateStrategy = 1
	GroupUpdateStrategy_OWNER_MANAGER        GroupUpdateStrategy = 2
	GroupUpdateStrategy_OWNER                GroupUpdateStrategy = 3
)
