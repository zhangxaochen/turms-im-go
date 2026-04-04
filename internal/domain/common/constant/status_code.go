package constant

// ResponseStatusCode defines the business response codes for Turms.
type ResponseStatusCode int32

const (
	ResponseStatusCode_OK ResponseStatusCode = 1000

	// Client Error
	ResponseStatusCode_INVALID_REQUEST              ResponseStatusCode = 1100
	ResponseStatusCode_CLIENT_REQUESTS_TOO_FREQUENT ResponseStatusCode = 1101
	ResponseStatusCode_ILLEGAL_ARGUMENT             ResponseStatusCode = 1102
	ResponseStatusCode_UNAUTHORIZED_REQUEST         ResponseStatusCode = 1106

	// Server Error
	ResponseStatusCode_SERVER_INTERNAL_ERROR ResponseStatusCode = 1200
	ResponseStatusCode_SERVER_UNAVAILABLE    ResponseStatusCode = 1201

	// Session
	ResponseStatusCode_CREATE_EXISTING_SESSION     ResponseStatusCode = 2001 // Add other session codes as requested
	ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED ResponseStatusCode = 2000
	ResponseStatusCode_LOGGING_IN_USER_NOT_ACTIVE  ResponseStatusCode = 2002

	ResponseStatusCode_LOGIN_TIMEOUT                           ResponseStatusCode = 2010
	ResponseStatusCode_UPDATE_HEARTBEAT_OF_NONEXISTENT_SESSION ResponseStatusCode = 2011

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
	ResponseStatusCode_QUERY_GROUP_JOIN_QUESTIONS_IS_DISABLED              ResponseStatusCode = 3800
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_CREATE_GROUP_QUESTION ResponseStatusCode = 3801
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_DELETE_GROUP_QUESTION ResponseStatusCode = 3802
	ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_QUESTION ResponseStatusCode = 3803
	ResponseStatusCode_CHECKING_GROUP_JOIN_QUESTION_ANSWER_IS_DISABLED     ResponseStatusCode = 3804
	ResponseStatusCode_ANSWER_GROUP_QUESTION_IS_DISABLED                   ResponseStatusCode = 3805
	ResponseStatusCode_ANY_GROUP_QUESTION_ANSWER_IS_INCORRECT              ResponseStatusCode = 3806

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
)

// IsServerError returns true if the turms business code implies a server-side runtime error.
func IsServerError(code int32) bool {
	// Roughly based on original Turms ResponseStatusCode mapping ranges:
	// Usually 1200 / 12xx maps to internal server errors.
	return code >= 1200 && code < 1300
}
