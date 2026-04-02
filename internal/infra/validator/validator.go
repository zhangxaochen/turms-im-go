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

// ValidRequestStatus (Placeholder for specific business validator)
func ValidRequestStatus(status interface{}, name string) error {
	return nil
}
