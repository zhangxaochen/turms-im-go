package validator

import (
	"fmt"
	"reflect"
	"time"

	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/codes"
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

// ValidRequestStatus
// @MappedFrom validRequestStatus(RequestStatus status)
func ValidRequestStatus(status interface{}, name string) error {
	// Simple stub for missing PB logic
	return nil
}

// @MappedFrom validResponseAction(ResponseAction action)
func ValidResponseAction(action interface{}) error {
	return nil
}

// @MappedFrom validDeviceType(DeviceType deviceType)
func ValidDeviceType(deviceType interface{}) error {
	return nil
}

// @MappedFrom validProfileAccess(ProfileAccessStrategy value)
func ValidProfileAccess(value interface{}) error {
	return nil
}

// @MappedFrom validRelationshipKey(UserRelationship.Key key)
func ValidRelationshipKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "UserRelationship key must not be null")
	}
	return nil
}

// @MappedFrom validRelationshipGroupKey(UserRelationshipGroup.Key key)
func ValidRelationshipGroupKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "UserRelationshipGroup key must not be null")
	}
	return nil
}

// @MappedFrom validGroupMemberKey(GroupMember.Key key)
func ValidGroupMemberKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupMember key must not be null")
	}
	return nil
}

// @MappedFrom validGroupMemberRole(GroupMemberRole role)
func ValidGroupMemberRole(role interface{}) error {
	return nil
}

// @MappedFrom validGroupBlockedUserKey(GroupBlockedUser.Key key)
func ValidGroupBlockedUserKey(key interface{}) error {
	if key == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "GroupBlockedUser key must not be null")
	}
	return nil
}

// @MappedFrom validNewGroupQuestion(NewGroupQuestion question)
func ValidNewGroupQuestion(question interface{}) error {
	if question == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "NewGroupQuestion must not be null")
	}
	return nil
}

// @MappedFrom validGroupQuestionIdAndAnswer(Map.Entry<Long, String> questionIdAndAnswer)
func ValidGroupQuestionIdAndAnswer(questionIdAndAnswer interface{}) error {
	if questionIdAndAnswer == nil {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "Map entry must not be null")
	}
	return nil
}
