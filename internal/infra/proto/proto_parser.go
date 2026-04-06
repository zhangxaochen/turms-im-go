package proto

import (
	"fmt"

	stdproto "google.golang.org/protobuf/proto"

	"im.turms/server/pkg/protocol"
)

// KindCase maps to TurmsRequest.KindCase in Java.
// It represents the field number of the oneof "kind" field in TurmsRequest.
// KIND_NOT_SET = 0 (no oneof set).
type KindCase int32

const (
	KindNotSet KindCase = 0
)

// kindCaseNames maps KindCase (proto field number) to the Java TurmsRequest.KindCase enum name.
// This mirrors Java's KindCase.name() for logging parity.
var kindCaseNames = map[KindCase]string{
	KindNotSet: "KIND_NOT_SET",
	3:    "CREATE_SESSION_REQUEST",
	4:    "DELETE_SESSION_REQUEST",
	5:    "QUERY_CONVERSATIONS_REQUEST",
	6:    "UPDATE_CONVERSATION_REQUEST",
	7:    "UPDATE_TYPING_STATUS_REQUEST",
	8:    "CREATE_MESSAGE_REQUEST",
	9:    "QUERY_MESSAGES_REQUEST",
	10:   "UPDATE_MESSAGE_REQUEST",
	11:   "CREATE_GROUP_MEMBERS_REQUEST",
	12:   "DELETE_GROUP_MEMBERS_REQUEST",
	13:   "QUERY_GROUP_MEMBERS_REQUEST",
	14:   "UPDATE_GROUP_MEMBER_REQUEST",
	100:  "QUERY_USER_PROFILES_REQUEST",
	101:  "QUERY_NEARBY_USERS_REQUEST",
	102:  "QUERY_USER_ONLINE_STATUSES_REQUEST",
	103:  "UPDATE_USER_LOCATION_REQUEST",
	104:  "UPDATE_USER_ONLINE_STATUS_REQUEST",
	105:  "UPDATE_USER_REQUEST",
	200:  "CREATE_FRIEND_REQUEST_REQUEST",
	201:  "CREATE_RELATIONSHIP_GROUP_REQUEST",
	202:  "CREATE_RELATIONSHIP_REQUEST",
	203:  "DELETE_FRIEND_REQUEST_REQUEST",
	204:  "DELETE_RELATIONSHIP_GROUP_REQUEST",
	205:  "DELETE_RELATIONSHIP_REQUEST",
	206:  "QUERY_FRIEND_REQUESTS_REQUEST",
	207:  "QUERY_RELATED_USER_IDS_REQUEST",
	208:  "QUERY_RELATIONSHIP_GROUPS_REQUEST",
	209:  "QUERY_RELATIONSHIPS_REQUEST",
	210:  "UPDATE_FRIEND_REQUEST_REQUEST",
	211:  "UPDATE_RELATIONSHIP_GROUP_REQUEST",
	212:  "UPDATE_RELATIONSHIP_REQUEST",
	300:  "CREATE_GROUP_REQUEST",
	301:  "DELETE_GROUP_REQUEST",
	302:  "QUERY_GROUPS_REQUEST",
	303:  "QUERY_JOINED_GROUP_IDS_REQUEST",
	304:  "QUERY_JOINED_GROUP_INFOS_REQUEST",
	305:  "UPDATE_GROUP_REQUEST",
	400:  "CREATE_GROUP_BLOCKED_USER_REQUEST",
	401:  "DELETE_GROUP_BLOCKED_USER_REQUEST",
	402:  "QUERY_GROUP_BLOCKED_USER_IDS_REQUEST",
	403:  "QUERY_GROUP_BLOCKED_USER_INFOS_REQUEST",
	500:  "CHECK_GROUP_JOIN_QUESTIONS_ANSWERS_REQUEST",
	501:  "CREATE_GROUP_INVITATION_REQUEST",
	502:  "CREATE_GROUP_JOIN_REQUEST_REQUEST",
	503:  "CREATE_GROUP_JOIN_QUESTIONS_REQUEST",
	504:  "DELETE_GROUP_INVITATION_REQUEST",
	505:  "DELETE_GROUP_JOIN_REQUEST_REQUEST",
	506:  "DELETE_GROUP_JOIN_QUESTIONS_REQUEST",
	507:  "QUERY_GROUP_INVITATIONS_REQUEST",
	508:  "QUERY_GROUP_JOIN_REQUESTS_REQUEST",
	509:  "QUERY_GROUP_JOIN_QUESTIONS_REQUEST",
	510:  "UPDATE_GROUP_INVITATION_REQUEST",
	511:  "UPDATE_GROUP_JOIN_QUESTION_REQUEST",
	512:  "UPDATE_GROUP_JOIN_REQUEST_REQUEST",
	900:  "CREATE_MEETING_REQUEST",
	901:  "DELETE_MEETING_REQUEST",
	902:  "QUERY_MEETINGS_REQUEST",
	903:  "UPDATE_MEETING_REQUEST",
	904:  "UPDATE_MEETING_INVITATION_REQUEST",
	1000: "DELETE_RESOURCE_REQUEST",
	1001: "QUERY_RESOURCE_DOWNLOAD_INFO_REQUEST",
	1002: "QUERY_RESOURCE_UPLOAD_INFO_REQUEST",
	1003: "QUERY_MESSAGE_ATTACHMENT_INFOS_REQUEST",
	1004: "UPDATE_MESSAGE_ATTACHMENT_INFO_REQUEST",
	1100: "DELETE_CONVERSATION_SETTINGS_REQUEST",
	1101: "QUERY_CONVERSATION_SETTINGS_REQUEST",
	1102: "UPDATE_CONVERSATION_SETTINGS_REQUEST",
}

// KindCaseName returns the Java-style enum name for a KindCase value.
// Returns "KIND_NOT_SET" for 0, and the enum name for known field numbers.
// Falls back to the numeric string for unknown field numbers.
func KindCaseName(kc KindCase) string {
	if name, ok := kindCaseNames[kc]; ok {
		return name
	}
	return fmt.Sprintf("%d", int32(kc))
}

// SimpleTurmsNotification maps to SimpleTurmsNotification record in Java.
// @MappedFrom SimpleTurmsNotification
type SimpleTurmsNotification struct {
	RequesterID int64
	CloseStatus *int32
	// RelayedRequestType is the KindCase (field number) of the relayed TurmsRequest.
	// Using KindCase (int32) instead of `any` for type safety.
	RelayedRequestType KindCase
}

func NewSimpleTurmsNotification(requesterID int64, closeStatus *int32, relayedRequestType KindCase) *SimpleTurmsNotification {
	return &SimpleTurmsNotification{
		RequesterID:        requesterID,
		CloseStatus:        closeStatus,
		RelayedRequestType: relayedRequestType,
	}
}

// SimpleTurmsRequest maps to SimpleTurmsRequest record in Java.
// @MappedFrom SimpleTurmsRequest
type SimpleTurmsRequest struct {
	RequestID int64
	// Type is the KindCase of the request.
	Type                 KindCase
	CreateSessionRequest *protocol.CreateSessionRequest
}

func NewSimpleTurmsRequest(requestID int64, reqType KindCase, createSessionReq *protocol.CreateSessionRequest) *SimpleTurmsRequest {
	return &SimpleTurmsRequest{
		RequestID:            requestID,
		Type:                 reqType,
		CreateSessionRequest: createSessionReq,
	}
}

// @MappedFrom toString()
func (r *SimpleTurmsRequest) ToString() string {
	csrStr := ToLogString(r.CreateSessionRequest)
	return fmt.Sprintf("SimpleTurmsRequest[requestId=%d, type=%v, createSessionRequest=%s]",
		r.RequestID, r.Type, csrStr)
}

// --- Protobuf wire format constants ---

// Wire types
const (
	wireTypeVarint       = 0
	wireTypeLenDelimited = 2
)

// TurmsNotification field numbers (from the .proto definition)
// requester_id = field 10, wire type 0 → tag = (10 << 3) | 0 = 80
// close_status = field 11, wire type 0 → tag = (11 << 3) | 0 = 88
// relayed_request = field 12, wire type 2 → tag = (12 << 3) | 2 = 98
const (
	notificationRequesterIDTag = 80
	notificationCloseStatusTag = 88
	notificationRelayedReqTag  = 98
)

// TurmsRequest field numbers
// request_id = field 1, wire type 0 → tag = 8
const (
	requestIDTag = 8
)

// undefinedRequestID matches Java's Long.MIN_VALUE sentinel
const undefinedRequestID = int64(-9223372036854775808)

// ResponseException is a simple error type for request parsing errors.
type ResponseException struct {
	Code    int
	Message string
}

func (e *ResponseException) Error() string {
	return fmt.Sprintf("ResponseException(code=%d): %s", e.Code, e.Message)
}

const (
	statusCodeIllegalArgument = 1400 // matches Java ResponseStatusCode.ILLEGAL_ARGUMENT
)

func newIllegalArgErr(msg string) *ResponseException {
	return &ResponseException{Code: statusCodeIllegalArgument, Message: msg}
}

// --- Varint decoding helpers ---

// readVarint reads a protobuf varint from data starting at pos.
// Returns the decoded value and the new position.
func readVarint(data []byte, pos int) (uint64, int, error) {
	var result uint64
	var shift uint
	for {
		if pos >= len(data) {
			return 0, pos, fmt.Errorf("unexpected end of data while reading varint")
		}
		b := data[pos]
		pos++
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, pos, nil
		}
		shift += 7
		if shift >= 64 {
			return 0, pos, fmt.Errorf("varint overflow")
		}
	}
}

// readInt64 reads a varint and interprets it as int64.
func readInt64(data []byte, pos int) (int64, int, error) {
	v, newPos, err := readVarint(data, pos)
	return int64(v), newPos, err
}

// readInt32 reads a varint and interprets it as int32.
func readInt32(data []byte, pos int) (int32, int, error) {
	v, newPos, err := readVarint(data, pos)
	return int32(v), newPos, err
}

// TurmsNotificationParser maps to TurmsNotificationParser in Java.
// @MappedFrom TurmsNotificationParser
type TurmsNotificationParser struct{}

// @MappedFrom parseSimpleNotification(CodedInputStream turmsRequestInputStream)
// ParseSimpleNotification parses a TurmsNotification protobuf wire-format byte slice into
// a SimpleTurmsNotification containing requesterId, optional closeStatus, and relayedRequestType.
func (p *TurmsNotificationParser) ParseSimpleNotification(data []byte) (*SimpleTurmsNotification, error) {
	if data == nil {
		return nil, newIllegalArgErr("The input stream must not be null")
	}

	var requesterID int64 = undefinedRequestID
	var requesterIDSet bool
	var closeStatus *int32
	var kindCase KindCase
	var kindSet bool
	done := false

	pos := 0
	for !done && pos < len(data) {
		tag, newPos, err := readVarint(data, pos)
		if err != nil {
			return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
		}
		pos = newPos
		if tag == 0 {
			break // end of message
		}

		switch int(tag) {
		case notificationRequesterIDTag: // field 10, varint
			if requesterIDSet {
				return nil, newIllegalArgErr("Not a valid notification: Duplicate requester ID")
			}
			v, np, err := readInt64(data, pos)
			if err != nil {
				return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
			}
			requesterID = v
			requesterIDSet = true
			pos = np

		case notificationCloseStatusTag: // field 11, varint
			if closeStatus != nil {
				return nil, newIllegalArgErr("Not a valid notification: Duplicate close status")
			}
			v, np, err := readInt32(data, pos)
			if err != nil {
				return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
			}
			closeStatus = &v
			pos = np

		case notificationRelayedReqTag: // field 12, length-delimited
			// Read the length prefix of the embedded TurmsRequest
			_, np, err := readVarint(data, pos)
			if err != nil {
				return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
			}
			pos = np
			// Read the inner tag to get the field number (= KindCase)
			if pos >= len(data) {
				break
			}
			innerTag, np2, err := readVarint(data, pos)
			if err != nil {
				return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
			}
			pos = np2
			kindFieldNumber := int32(innerTag >> 3)
			kindCase = KindCase(kindFieldNumber)
			kindSet = true
			done = true

		default:
			// Skip unknown field based on wire type
			wireType := tag & 0x7
			switch wireType {
			case wireTypeVarint:
				_, np, err := readVarint(data, pos)
				if err != nil {
					return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
				}
				pos = np
			case wireTypeLenDelimited:
				length, np, err := readVarint(data, pos)
				if err != nil {
					return nil, newIllegalArgErr("Not a valid notification: " + err.Error())
				}
				endPos := np + int(length)
				if endPos > len(data) {
					return nil, newIllegalArgErr("Not a valid notification: truncated length-delimited field")
				}
				pos = endPos
			case 1: // 64-bit
				if pos+8 > len(data) {
					return nil, newIllegalArgErr("Not a valid notification: truncated 64-bit field")
				}
				pos += 8
			case 5: // 32-bit
				if pos+4 > len(data) {
					return nil, newIllegalArgErr("Not a valid notification: truncated 32-bit field")
				}
				pos += 4
			default:
				return nil, newIllegalArgErr(fmt.Sprintf("Not a valid notification: unknown wire type %d", wireType))
			}
		}
	}

	if !requesterIDSet {
		return nil, newIllegalArgErr("Not a valid notification: No requester ID")
	}
	if !kindSet || kindCase == KindNotSet {
		return nil, newIllegalArgErr(fmt.Sprintf("Not a valid notification: Unknown request type: %d", kindCase))
	}

	return NewSimpleTurmsNotification(requesterID, closeStatus, kindCase), nil
}

// TurmsRequestParser maps to TurmsRequestParser in Java.
// @MappedFrom TurmsRequestParser
type TurmsRequestParser struct{}

// @MappedFrom parseSimpleRequest(CodedInputStream turmsRequestInputStream)
// ParseSimpleRequest parses a TurmsRequest protobuf wire-format byte slice into a SimpleTurmsRequest.
func (p *TurmsRequestParser) ParseSimpleRequest(data []byte) (*SimpleTurmsRequest, error) {
	if data == nil {
		return nil, newIllegalArgErr("The input stream must not be null")
	}

	requestID := undefinedRequestID
	var requestIDSet bool
	var kindCase KindCase
	var kindPos int // position after the inner kind tag for CreateSessionRequest parsing

	pos := 0
	for pos < len(data) {
		tag, newPos, err := readVarint(data, pos)
		if err != nil {
			return nil, newIllegalArgErr("Not a valid request: " + err.Error())
		}
		pos = newPos
		if tag == 0 {
			break
		}

		if int(tag) == requestIDTag { // field 1, varint
			if requestIDSet {
				return nil, newIllegalArgErr("Not a valid request: Duplicate request ID")
			}
			v, np, err := readInt64(data, pos)
			if err != nil {
				return nil, newIllegalArgErr("Not a valid request: " + err.Error())
			}
			requestID = v
			requestIDSet = true
			pos = np
		} else {
			// This tag's field number is the KindCase
			kindFieldNumber := int32(tag >> 3)
			kindCase = KindCase(kindFieldNumber)
			kindPos = pos
			break
		}
	}

	if !requestIDSet {
		return nil, newIllegalArgErr("Not a valid request: No request ID")
	}
	if kindCase == KindNotSet {
		return nil, newIllegalArgErr(fmt.Sprintf("Not a valid request: Unknown request type: %d", kindCase))
	}

	// CREATE_SESSION_REQUEST has KindCase value matching proto field number for create_session_request.
	// Field create_session_request = 3 in TurmsRequest.proto → KindCase = 3.
	const createSessionRequestKind KindCase = 3

	var createSessionReq *protocol.CreateSessionRequest
	if kindCase == createSessionRequestKind {
		// The oneof field is length-delimited; read length then unmarshal
		length, np, err := readVarint(data, kindPos)
		if err != nil {
			return nil, newIllegalArgErr("Not a valid request: " + err.Error())
		}
		subdata := data[np : np+int(length)]
		var req protocol.CreateSessionRequest
		if err := stdproto.Unmarshal(subdata, &req); err != nil {
			return nil, newIllegalArgErr("Not a valid request: " + err.Error())
		}
		createSessionReq = &req
	}

	return NewSimpleTurmsRequest(requestID, kindCase, createSessionReq), nil
}

// unmarshalProto uses google.golang.org/protobuf/proto.Unmarshal to decode wire-format bytes.
func unmarshalProto(data []byte, msg stdproto.Message) error {
	return stdproto.Unmarshal(data, msg)
}
