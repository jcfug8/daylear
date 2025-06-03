package namer

// ErrInvalidParent is returned when a parent resource name is invalid or does not match any known pattern.
type ErrInvalidParent struct {
	msg string
}

// Error implements the error interface for ErrInvalidParent.
func (e ErrInvalidParent) Error() string {
	return e.msg
}

// ErrInvalidName is returned when a resource name is invalid or does not match any known pattern.
type ErrInvalidName struct {
	msg string
}

// Error implements the error interface for ErrInvalidName.
func (e ErrInvalidName) Error() string {
	return e.msg
}

// ErrInvalidPatternIndex is returned when a pattern index is out of range or invalid.
type ErrInvalidPatternIndex struct {
	msg string
}

// Error implements the error interface for ErrInvalidPatternIndex.
func (e ErrInvalidPatternIndex) Error() string {
	return e.msg
}

// ErrNoPatternFound is returned when no suitable resource name pattern can be found for a given input.
type ErrNoPatternFound struct {
	msg string
}

// Error implements the error interface for ErrNoPatternFound.
func (e ErrNoPatternFound) Error() string {
	return e.msg
}

// ErrInvalidField is returned when a required struct field is missing or does not match the expected pattern key.
type ErrInvalidField struct {
	msg string
}

// Error implements the error interface for ErrInvalidField.
func (e ErrInvalidField) Error() string {
	return e.msg
}

// ErrInvalidFieldValue is returned when a struct field value is invalid (e.g., zero, empty, or nil) for resource name formatting/parsing.
type ErrInvalidFieldValue struct {
	msg string
}

// Error implements the error interface for ErrInvalidFieldValue.
func (e ErrInvalidFieldValue) Error() string {
	return e.msg
}

// ErrInvalidFieldType is returned when a struct field is of an unsupported type for resource name formatting/parsing.
type ErrInvalidFieldType struct {
	msg string
}

// Error implements the error interface for ErrInvalidFieldType.
func (e ErrInvalidFieldType) Error() string {
	return e.msg
}
