package exception

import (
	"fmt"
)

type TurmsError struct {
	Code    int32
	Message string
}

func (e *TurmsError) Error() string {
	return fmt.Sprintf("TurmsError(code=%d, message=%s)", e.Code, e.Message)
}

func NewTurmsError(code int32, message string) *TurmsError {
	return &TurmsError{
		Code:    code,
		Message: message,
	}
}

// Get returns a TurmsError with the given code and a default message (currently just returning the code string)
// In a full implementation, this could map codes to localized strings.
func Get(code int32) *TurmsError {
	return &TurmsError{
		Code:    code,
		Message: fmt.Sprintf("Error code: %d", code),
	}
}

// IsCode checks if an error is a TurmsError with a specific code
func IsCode(err error, code int32) bool {
	if te, ok := err.(*TurmsError); ok {
		return te.Code == code
	}
	return false
}

// IsDuplicateKeyError returns true if the error is a MongoDB duplicate key error (code 11000 or 11001)
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// MongoDB error code 11000 is for duplicate key
	const mongoDuplicateKeyCode = 11000
	const mongoDuplicateKeyCodeWriteConflict = 11001

	// In the Go MongoDB driver, error codes can be checked via CommandError or WriteException
	type mongoError interface {
		HasErrorCode(code int) bool
	}
	if me, ok := err.(mongoError); ok {
		return me.HasErrorCode(mongoDuplicateKeyCode) || me.HasErrorCode(mongoDuplicateKeyCodeWriteConflict)
	}
	return false
}
