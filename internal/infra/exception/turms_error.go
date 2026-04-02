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
