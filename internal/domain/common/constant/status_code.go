package constant

import "fmt"

// ResponseStatusCode defines the business response codes for Turms.
type ResponseStatusCode int32

const (
	ResponseStatusCode_OK ResponseStatusCode = 1000

	// Client Error
	ResponseStatusCode_INVALID_REQUEST              ResponseStatusCode = 1100
	ResponseStatusCode_CLIENT_REQUESTS_TOO_FREQUENT ResponseStatusCode = 1101
	ResponseStatusCode_ILLEGAL_ARGUMENT             ResponseStatusCode = 1102
	ResponseStatusCode_UNAUTHORIZED_REQUEST         ResponseStatusCode = 1106
	ResponseStatusCode_NO_CONTENT                   ResponseStatusCode = 1001
	ResponseStatusCode_UNSUPPORTED_CLIENT_VERSION   ResponseStatusCode = 103

	// Server Error
	ResponseStatusCode_SERVER_INTERNAL_ERROR ResponseStatusCode = 1200
	ResponseStatusCode_SERVER_UNAVAILABLE    ResponseStatusCode = 1201

	// Session
	ResponseStatusCode_CREATE_EXISTING_SESSION                ResponseStatusCode = 2001 // Add other session codes as requested
	ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED            ResponseStatusCode = 2000
	ResponseStatusCode_LOGGING_IN_USER_NOT_ACTIVE             ResponseStatusCode = 2002
	ResponseStatusCode_SESSION_SIMULTANEOUS_CONFLICTS_DECLINE ResponseStatusCode = 2004
	ResponseStatusCode_LOGIN_FROM_FORBIDDEN_DEVICE_TYPE       ResponseStatusCode = 1103

	ResponseStatusCode_LOGIN_TIMEOUT                           ResponseStatusCode = 2010
	ResponseStatusCode_UPDATE_HEARTBEAT_OF_NONEXISTENT_SESSION ResponseStatusCode = 2011
	ResponseStatusCode_UPDATE_NON_EXISTING_SESSION_STATUS      ResponseStatusCode = 2011

	// Group Error
	ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP                ResponseStatusCode = 3000
	ResponseStatusCode_NOT_GROUP_OWNER_TO_UPDATE_GROUP_INFO            ResponseStatusCode = 3001
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_INFO ResponseStatusCode = 3002
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_UPDATE_GROUP_INFO           ResponseStatusCode = 3003

	ResponseStatusCode_NOT_GROUP_OWNER_TO_TRANSFER_GROUP ResponseStatusCode = 3201
	ResponseStatusCode_NOT_GROUP_OWNER_TO_DELETE_GROUP   ResponseStatusCode = 3202
	ResponseStatusCode_GROUP_SUCCESSOR_NOT_GROUP_MEMBER  ResponseStatusCode = 3203
	ResponseStatusCode_MAX_OWNED_GROUPS_REACHED          ResponseStatusCode = 3205
	ResponseStatusCode_TRANSFER_NONEXISTENT_GROUP        ResponseStatusCode = 3206

	ResponseStatusCode_ADD_USER_TO_GROUP_REQUIRING_USERS_APPROVAL                    ResponseStatusCode = 3400
	ResponseStatusCode_ADD_USER_TO_INACTIVE_GROUP                                    ResponseStatusCode = 3401
	ResponseStatusCode_NOT_GROUP_OWNER_TO_ADD_GROUP_MANAGER                          ResponseStatusCode = 3402
	ResponseStatusCode_ADD_USER_TO_GROUP_WITH_SIZE_LIMIT_REACHED                     ResponseStatusCode = 3403
	ResponseStatusCode_ADD_BLOCKED_USER_TO_GROUP                                     ResponseStatusCode = 3404
	ResponseStatusCode_NOT_GROUP_OWNER_TO_ADD_GROUP_MEMBER                           ResponseStatusCode = 3405
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_GROUP_MEMBER                ResponseStatusCode = 3406
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_ADD_GROUP_MEMBER                          ResponseStatusCode = 3407
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_GROUP_MEMBER             ResponseStatusCode = 3408
	ResponseStatusCode_NOT_GROUP_OWNER_TO_REMOVE_GROUP_OWNER_OR_MANAGER              ResponseStatusCode = 3409
	ResponseStatusCode_NOT_GROUP_OWNER_TO_UPDATE_GROUP_MEMBER_ROLE                   ResponseStatusCode = 3410
	ResponseStatusCode_UPDATE_GROUP_MEMBER_ROLE_OF_NONEXISTENT_GROUP                 ResponseStatusCode = 3411
	ResponseStatusCode_NOT_GROUP_OWNER_TO_UPDATE_GROUP_MEMBER_INFO                   ResponseStatusCode = 3412
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_MEMBER_INFO        ResponseStatusCode = 3413
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_UPDATE_GROUP_MEMBER_INFO                  ResponseStatusCode = 3414
	ResponseStatusCode_UPDATE_GROUP_MEMBER_INFO_OF_NONEXISTENT_GROUP                 ResponseStatusCode = 3415
	ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP_MEMBER                       ResponseStatusCode = 3416
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_MUTE_GROUP_MEMBER               ResponseStatusCode = 3417
	ResponseStatusCode_MUTE_GROUP_MEMBER_WITH_ROLE_EQUAL_TO_OR_HIGHER_THAN_REQUESTER ResponseStatusCode = 3418
	ResponseStatusCode_MUTE_GROUP_MEMBER_OF_NONEXISTENT_GROUP                        ResponseStatusCode = 3419
	ResponseStatusCode_MUTE_NONEXISTENT_GROUP_MEMBER                                 ResponseStatusCode = 3420
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_QUERY_GROUP_MEMBER_INFO                   ResponseStatusCode = 3421
	ResponseStatusCode_USER_JOIN_GROUP_WITHOUT_ACCEPTING_GROUP_INVITATION            ResponseStatusCode = 3422
	ResponseStatusCode_USER_JOIN_GROUP_WITHOUT_ANSWERING_GROUP_QUESTION              ResponseStatusCode = 3423
	ResponseStatusCode_USER_JOIN_GROUP_WITHOUT_SENDING_GROUP_JOIN_REQUEST            ResponseStatusCode = 3424

	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_BLOCKED_USER    ResponseStatusCode = 3500
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_BLOCKED_USER ResponseStatusCode = 3501

	// Group - Join Request
	ResponseStatusCode_GROUP_JOIN_REQUEST_IS_DISABLED                          ResponseStatusCode = 3600
	ResponseStatusCode_ADD_USER_TO_GROUP_REQUIRING_INVITATION                  ResponseStatusCode = 3601
	ResponseStatusCode_RECALLING_GROUP_JOIN_REQUEST_IS_DISABLED                ResponseStatusCode = 3602
	ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST                   ResponseStatusCode = 3603
	ResponseStatusCode_NOT_SENDER_TO_RECALL_GROUP_JOIN_REQUEST                 ResponseStatusCode = 3604
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_QUERY_GROUP_JOIN_REQUEST  ResponseStatusCode = 3605
	ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST                   ResponseStatusCode = 3606
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_JOIN_REQUEST ResponseStatusCode = 3607

	// Group - Invitation
	ResponseStatusCode_SENDING_GROUP_INVITATION_IS_DISABLED                            ResponseStatusCode = 3700
	ResponseStatusCode_QUERY_GROUP_INVITATIONS_IS_DISABLED                             ResponseStatusCode = 3701
	ResponseStatusCode_NOT_GROUP_OWNER_TO_SEND_GROUP_INVITATION                        ResponseStatusCode = 3702
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_SEND_GROUP_INVITATION             ResponseStatusCode = 3703
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_SEND_GROUP_INVITATION                       ResponseStatusCode = 3704
	ResponseStatusCode_RECALLING_GROUP_INVITATION_IS_DISABLED                          ResponseStatusCode = 3705
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_RECALL_GROUP_INVITATION           ResponseStatusCode = 3706
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_OR_SENDER_TO_RECALL_GROUP_INVITATION ResponseStatusCode = 3707
	ResponseStatusCode_RECALL_NON_PENDING_GROUP_INVITATION                             ResponseStatusCode = 3708
	ResponseStatusCode_UPDATE_NON_PENDING_GROUP_INVITATION                             ResponseStatusCode = 3709
	ResponseStatusCode_NOT_INVITEE_TO_UPDATE_GROUP_INVITATION                          ResponseStatusCode = 3710
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_QUERY_GROUP_INVITATION            ResponseStatusCode = 3711

	// Group - Question
	ResponseStatusCode_QUERY_GROUP_JOIN_QUESTIONS_IS_DISABLED                ResponseStatusCode = 3800
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_CREATE_GROUP_QUESTION   ResponseStatusCode = 3801
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_DELETE_GROUP_QUESTION   ResponseStatusCode = 3802
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_QUESTION   ResponseStatusCode = 3803
	ResponseStatusCode_CHECKING_GROUP_JOIN_QUESTION_ANSWER_IS_DISABLED       ResponseStatusCode = 3804
	ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED                     ResponseStatusCode = 3805
	ResponseStatusCode_ANY_GROUP_QUESTION_ANSWER_IS_INCORRECT                ResponseStatusCode = 3806
	ResponseStatusCode_GROUP_QUESTION_ANSWERER_HAS_BEEN_BLOCKED              ResponseStatusCode = 3807
	ResponseStatusCode_GROUP_MEMBER_ANSWER_GROUP_QUESTION                    ResponseStatusCode = 3808
	ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_INACTIVE_GROUP              ResponseStatusCode = 3809
	ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_GROUP_USING_JOIN_REQUEST    ResponseStatusCode = 3810
	ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_GROUP_USING_INVITATION      ResponseStatusCode = 3811
	ResponseStatusCode_CREATE_GROUP_QUESTION_FOR_GROUP_USING_MEMBERSHIP_REQUEST ResponseStatusCode = 3812
	ResponseStatusCode_ANSWER_GROUP_QUESTION_OF_INACTIVE_GROUP              ResponseStatusCode = 3813
	ResponseStatusCode_GROUP_INVITER_NOT_MEMBER                              ResponseStatusCode = 3814

	ResponseStatusCode_BLOCKED_USER_SEND_GROUP_JOIN_REQUEST ResponseStatusCode = 3404
	ResponseStatusCode_USER_ALREADY_GROUP_MEMBER            ResponseStatusCode = 3405
	ResponseStatusCode_GROUP_JOIN_STRATEGY_NOT_JOIN_REQUEST ResponseStatusCode = 3407

	ResponseStatusCode_SEND_GROUP_INVITATION_TO_GROUP_MEMBER ResponseStatusCode = 3712
	ResponseStatusCode_SEND_GROUP_INVITATION_TO_BLOCKED_USER ResponseStatusCode = 3713

	ResponseStatusCode_SEND_MESSAGE_TO_INACTIVE_GROUP            ResponseStatusCode = 5004
	ResponseStatusCode_SEND_MESSAGE_TO_MUTED_GROUP               ResponseStatusCode = 5005
	ResponseStatusCode_SEND_MESSAGE_TO_NONEXISTENT_GROUP         ResponseStatusCode = 5006
	ResponseStatusCode_MUTED_GROUP_MEMBER_SEND_MESSAGE           ResponseStatusCode = 5008
	ResponseStatusCode_NOT_SPEAKABLE_GROUP_GUEST_TO_SEND_MESSAGE ResponseStatusCode = 5009
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_SEND_MESSAGE          ResponseStatusCode = 5010
	ResponseStatusCode_MUTED_MEMBER_SEND_MESSAGE                 ResponseStatusCode = 5011
	ResponseStatusCode_BLOCKED_USER_SEND_GROUP_MESSAGE           ResponseStatusCode = 5012
	ResponseStatusCode_CONFERENCE_NOT_IMPLEMENTED                ResponseStatusCode = 8000

	// Conversation Error
	ResponseStatusCode_UPDATING_READ_DATE_IS_DISABLED                             ResponseStatusCode = 5000
	ResponseStatusCode_UPDATING_READ_DATE_OF_NONEXISTENT_GROUP_CONVERSATION       ResponseStatusCode = 5001
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_UPDATE_READ_DATE_OF_GROUP_CONVERSATION ResponseStatusCode = 5002
	ResponseStatusCode_UPDATING_READ_DATE_IS_DISABLED_BY_GROUP                    ResponseStatusCode = 5003
	ResponseStatusCode_MOVING_READ_DATE_FORWARD_IS_DISABLED                       ResponseStatusCode = 5007
	ResponseStatusCode_UPDATING_TYPING_STATUS_IS_DISABLED                         ResponseStatusCode = 5100
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_SEND_TYPING_STATUS                     ResponseStatusCode = 5101
	ResponseStatusCode_NOT_FRIEND_TO_SEND_TYPING_STATUS                           ResponseStatusCode = 5102

	ResponseStatusCode_CREATE_MEETING_EXCEEDING_MAX_ACTIVE_MEETING_COUNT ResponseStatusCode = 8100
	ResponseStatusCode_NOT_CREATOR_TO_CANCEL_MEETING                     ResponseStatusCode = 8101
	ResponseStatusCode_CANCELING_MEETING_IS_DISABLED                     ResponseStatusCode = 8102
	ResponseStatusCode_CANCEL_NONEXISTENT_MEETING                        ResponseStatusCode = 8103
	ResponseStatusCode_NOT_CREATOR_TO_UPDATE_MEETING_PASSWORD            ResponseStatusCode = 8104
	ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_MEETING                ResponseStatusCode = 8105

	ResponseStatusCode_ACCEPT_MEETING_INVITATION_WITH_WRONG_PASSWORD ResponseStatusCode = 8200
	ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_CANCELED_MEETING ResponseStatusCode = 8201
	ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_ENDED_MEETING    ResponseStatusCode = 8202
	ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_EXPIRED_MEETING  ResponseStatusCode = 8203
	ResponseStatusCode_ACCEPT_MEETING_INVITATION_OF_PENDING_MEETING  ResponseStatusCode = 8204
	ResponseStatusCode_ACCEPT_NONEXISTENT_MEETING_INVITATION         ResponseStatusCode = 8205
	ResponseStatusCode_NOT_GROUP_MEMBER_TO_CREATE_MEETING            ResponseStatusCode = 5025 // Kept as added by previous devs
)

// Reason returns the default explanation for a specific status code.
// Note: In a full Go port, we would want to generate this from the proto or Java source.
func (c ResponseStatusCode) Reason() string {
	switch c {
	case ResponseStatusCode_OK:
		return "ok"
	case ResponseStatusCode_SERVER_INTERNAL_ERROR:
		return "Internal server error"
	case ResponseStatusCode_SERVER_UNAVAILABLE:
		return "The server is unavailable"
	case ResponseStatusCode_INVALID_REQUEST:
		return "The client request is invalid"
	case ResponseStatusCode_ILLEGAL_ARGUMENT:
		return "The arguments are illegal"
	case ResponseStatusCode_UNAUTHORIZED_REQUEST:
		return "The request is unauthorized"
	default:
		return fmt.Sprintf("status code %d", int(c))
	}
}

// IsServerError returns true if the turms business code implies a server-side runtime error.
func IsServerError(code int32) bool {
	// Roughly based on original Turms ResponseStatusCode mapping ranges:
	// Usually 1200 / 12xx maps to internal server errors.
	return code >= 1200 && code < 1300
}
