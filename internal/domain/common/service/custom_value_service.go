package service

// CustomValueService maps to CustomValueService.java
// It's responsible for parsing custom property values (like enums, ints, strings matching regexes, etc.)
// from Protobuf Values.
// @MappedFrom CustomValueService
type CustomValueService struct {
	exceptionMessagePrefixValueOf      string
	exceptionMessagePrefixString       string
	exceptionMessagePrefixStringLength string
	exceptionMessagePrefixArray        string
}

// NewCustomValueService creates a new CustomValueService (abstract in Java).
func NewCustomValueService(valueOf, str, strLen, arr string) *CustomValueService {
	return &CustomValueService{
		exceptionMessagePrefixValueOf:      valueOf,
		exceptionMessagePrefixString:       str,
		exceptionMessagePrefixStringLength: strLen,
		exceptionMessagePrefixArray:        arr,
	}
}

// TODO: Implement `parseValue` based on `CustomValueOneOfProperties`
