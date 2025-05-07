package filtering

import "fmt"

// NewOperatorError returns an error for an operator.
func NewOperatorError(format string, args ...any) error {
	return &operatorError{
		message: fmt.Sprintf(format, args...),
	}
}

type operatorError struct{ message string }

// Error returns the error message.
func (err *operatorError) Error() string {
	return err.message
}

// NewPathError returns an error for a path.
func NewPathError(format string, args ...any) error {
	return &pathError{
		message: fmt.Sprintf(format, args...),
	}
}

type pathError struct{ message string }

// Error returns the error message.
func (err *pathError) Error() string {
	return err.message
}

// NewValueError returns an error for a value.
func NewValueError(format string, args ...any) error {
	return &valueError{
		message: fmt.Sprintf(format, args...),
	}
}

type valueError struct{ message string }

// Error returns the error message.
func (err *valueError) Error() string {
	return err.message
}
