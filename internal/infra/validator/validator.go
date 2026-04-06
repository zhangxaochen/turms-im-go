package validator

import (
	"fmt"
	"reflect"
	"time"

	"im.turms/server/internal/domain/group/dto"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/codes"
	"im.turms/server/pkg/protocol"
)

// NotNull returns an IllegalArgument error if value is nil
func NotNull(value interface{}, name string) error {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), name+" must not be null")
	}
	return nil
}

// NotEmpty returns an IllegalArgument error if a slice or map is empty
func NotEmpty(value interface{}, name string) error {
	if value == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), name+" must not be empty")
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return exception.NewTurmsError(int32(codes.IllegalArgument), name+" must not be empty")
		}
	}
	return nil
}

// MaxLength returns an IllegalArgument error if string length exceeds max
func MaxLength(value *string, name string, max int) error {
	if value != nil && len(*value) > max {
		return exception.NewTurmsError(int32(codes.IllegalArgument), fmt.Sprintf("%s must not exceed %d characters", name, max))
	}
	return nil
}

// NotEquals returns an IllegalArgument error if v1 equals v2
func NotEquals(v1, v2 interface{}, message string) error {
	if reflect.DeepEqual(v1, v2) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), message)
	}
	return nil
}

// ShouldTrue returns an IllegalArgument error if condition is false
func ShouldTrue(condition bool, message string) error {
	if !condition {
		return exception.NewTurmsError(int32(codes.IllegalArgument), message)
	}
	return nil
}

// AreAllNull returns true if all values are nil
func AreAllNull(values ...interface{}) bool {
	for _, v := range values {
		if v != nil {
			rv := reflect.ValueOf(v)
			if rv.Kind() != reflect.Ptr || !rv.IsNil() {
				return false
			}
		}
	}
	return true
}

// PastOrPresent returns an IllegalArgument error if date is in the future
func PastOrPresent(date *time.Time, name string) error {
	if date != nil && date.After(time.Now().Add(1*time.Minute)) { // Allow 1 min drift
		return exception.NewTurmsError(int32(codes.IllegalArgument), name+" must be in the past or present")
	}
	return nil
}

// isValidEnum checks if a protobuf enum value is recognized (not UNRECOGNIZED).
// In Go protobuf, there is no UNRECOGNIZED constant; instead we check if the value
// exists in the enum's name map.
func isValidEnum(val int32, nameMap map[int32]string) bool {
	_, ok := nameMap[val]
	return ok
}

// ValidRequestStatus validates that the RequestStatus is recognized.
// Bug fix: Java checks status == RequestStatus.UNRECOGNIZED and throws ILLEGAL_ARGUMENT.
// Go protobuf doesn't have UNRECOGNIZED constant; check via enum map instead.
// @MappedFrom validRequestStatus(RequestStatus status)
func ValidRequestStatus(status protocol.RequestStatus, name string) error {
	if !isValidEnum(int32(status), protocol.RequestStatus_name) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), name+" must be a valid RequestStatus")
	}
	return nil
}

// ValidResponseAction validates that the ResponseAction is recognized.
// Bug fix: Java checks action == ResponseAction.UNRECOGNIZED and throws.
// @MappedFrom validResponseAction(ResponseAction action)
func ValidResponseAction(action protocol.ResponseAction) error {
	if !isValidEnum(int32(action), protocol.ResponseAction_name) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "ResponseAction must be a valid ResponseAction")
	}
	return nil
}

// ValidDeviceType validates that the DeviceType is recognized.
// Bug fix: Java checks deviceType == DeviceType.UNRECOGNIZED and throws.
// @MappedFrom validDeviceType(DeviceType deviceType)
func ValidDeviceType(deviceType protocol.DeviceType) error {
	if !isValidEnum(int32(deviceType), protocol.DeviceType_name) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "DeviceType must be a valid DeviceType")
	}
	return nil
}

// ValidProfileAccess validates that the ProfileAccessStrategy is recognized.
// Bug fix: Java checks value == ProfileAccessStrategy.UNRECOGNIZED and throws.
// @MappedFrom validProfileAccess(ProfileAccessStrategy value)
func ValidProfileAccess(value protocol.ProfileAccessStrategy) error {
	if !isValidEnum(int32(value), protocol.ProfileAccessStrategy_name) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "ProfileAccessStrategy must be a valid ProfileAccessStrategy")
	}
	return nil
}

// ValidRelationshipKey validates a UserRelationship.Key.
// Bug fix: Java also checks key.getOwnerId() == null and key.getRelatedUserId() == null.
// @MappedFrom validRelationshipKey(UserRelationship.Key key)
func ValidRelationshipKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "UserRelationship key must not be null")
	}
	// Check key fields via reflection for the struct with OwnerID and RelatedUserID
	rv := reflect.ValueOf(key)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "UserRelationship key must not be null")
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		ownerID := rv.FieldByName("OwnerID")
		relatedUserID := rv.FieldByName("RelatedUserID")
		if ownerID.IsValid() && ownerID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The owner ID in the user relationship key must not be null")
		}
		if relatedUserID.IsValid() && relatedUserID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The related user ID in the user relationship key must not be null")
		}
	}
	return nil
}

// ValidRelationshipGroupKey validates a UserRelationshipGroup.Key.
// Bug fix: Java also checks key.getOwnerId() == null and key.getGroupIndex() == null.
// @MappedFrom validRelationshipGroupKey(UserRelationshipGroup.Key key)
func ValidRelationshipGroupKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "UserRelationshipGroup key must not be null")
	}
	rv := reflect.ValueOf(key)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "UserRelationshipGroup key must not be null")
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		ownerID := rv.FieldByName("OwnerID")
		groupIndex := rv.FieldByName("Index")
		if ownerID.IsValid() && ownerID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The owner ID in the user relationship group key must not be null")
		}
		if groupIndex.IsValid() && groupIndex.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The group index in the user relationship group key must not be null")
		}
	}
	return nil
}

// ValidGroupMemberKey validates a GroupMember.Key.
// Bug fix: Java also checks key.getGroupId() == null and key.getUserId() == null.
// @MappedFrom validGroupMemberKey(GroupMember.Key key)
func ValidGroupMemberKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupMember key must not be null")
	}
	rv := reflect.ValueOf(key)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupMember key must not be null")
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		groupID := rv.FieldByName("GroupID")
		userID := rv.FieldByName("UserID")
		if groupID.IsValid() && groupID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The group ID in the group member key must not be null")
		}
		if userID.IsValid() && userID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The user ID in the group member key must not be null")
		}
	}
	return nil
}

// ValidGroupMemberRole validates that the GroupMemberRole is recognized.
// Bug fix: Java checks role == GroupMemberRole.UNRECOGNIZED and throws.
// @MappedFrom validGroupMemberRole(GroupMemberRole role)
func ValidGroupMemberRole(role protocol.GroupMemberRole) error {
	if !isValidEnum(int32(role), protocol.GroupMemberRole_name) {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupMemberRole must be a valid GroupMemberRole")
	}
	return nil
}

// ValidGroupBlockedUserKey validates a GroupBlockedUser.Key.
// Bug fix: Java also checks key.getGroupId() == null and key.getUserId() == null.
// @MappedFrom validGroupBlockedUserKey(GroupBlockedUser.Key key)
func ValidGroupBlockedUserKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupBlockedUser key must not be null")
	}
	rv := reflect.ValueOf(key)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupBlockedUser key must not be null")
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		groupID := rv.FieldByName("GroupID")
		userID := rv.FieldByName("UserID")
		if groupID.IsValid() && groupID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The group ID in the group blocked user key must not be null")
		}
		if userID.IsValid() && userID.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The user ID in the group blocked user key must not be null")
		}
	}
	return nil
}

// ValidNewGroupQuestion validates a NewGroupQuestion.
// Bug fix: Already checks empty answers and null/negative score, matching Java.
// @MappedFrom validNewGroupQuestion(NewGroupQuestion question)
func ValidNewGroupQuestion(question *dto.NewGroupQuestion) error {
	if question == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "NewGroupQuestion must not be null")
	}
	if len(question.Answers) == 0 {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "The answers of a new group question should not be empty")
	}
	if question.Score == nil || *question.Score < 0 {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "The score of a new group question should not be null and must be greater than or equal to 0")
	}
	return nil
}

// ValidGroupQuestionIdAndAnswer validates a questionIdAndAnswer entry.
// Bug fix: Java also checks key == null and value == null.
// @MappedFrom validGroupQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)
func ValidGroupQuestionIdAndAnswer(questionIdAndAnswer interface{}) error {
	if questionIdAndAnswer == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "Map entry must not be null")
	}
	// Check for struct with Key/Value fields (Go equivalent of Map.Entry)
	rv := reflect.ValueOf(questionIdAndAnswer)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "Map entry must not be null")
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		keyField := rv.FieldByName("Key")
		valueField := rv.FieldByName("Value")
		if keyField.IsValid() && keyField.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The question ID must not be null")
		}
		if valueField.IsValid() && valueField.IsZero() {
			return exception.NewTurmsError(int32(codes.IllegalArgument), "The answer must not be null")
		}
	}
	return nil
}
