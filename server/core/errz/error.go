package errz

import (
	"fmt"
	"strings"

	"github.com/jcfug8/daylear/server/core/mapz"

	"google.golang.org/grpc/status"
)

// Sanitize sanitizes the error message.
func Sanitize(toSanitize error) *Error {
	err := Wrap(toSanitize)
	msg := strings.SplitN(err.Msg, ":", 2)
	err.Msg = strings.TrimSpace(msg[0])
	return err
}

// Wrapf creates a new formatted core error from the error provided. If the provided error
// is already of type core error it is copied. Otherwise it returns a new core
// error using code Unknown.
func Wrapf(format string, args ...any) *Error {
	err := Wrap(fmt.Errorf(format, args...))

	for _, arg := range args {
		if arg, ok := arg.(*Error); ok {
			err.Code = arg.Code

			for k, v := range mapz.CopyMap(arg.Context) {
				err.Context[k] = v
			}
		}
	}

	return err
}

// Wrap creates a new core error from the error provided. If the provided error
// is already of type core error it is copied. Otherwise it returns a new core
// error using code Unknown.
func Wrap(err error) *Error {
	switch err := err.(type) {
	case *Error:
		return &Error{
			Context: mapz.CopyMap(err.Context),
			Code:    err.Code,
			Msg:     err.Msg,
		}
	case nil:
		return nil
	}

	return &Error{
		Code:    Unknown,
		Msg:     err.Error(),
		Context: make(mapz.Map),
	}
}

// NewAlreadyExists creates a new core error with code AlreadyExists and message
// formatted according to the provided format specifier.
func NewAlreadyExists(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(AlreadyExists)
}

// NewCancelled creates a new core error with code Cancelled and message
// formatted according to the provided format specifier.
func NewCancelled(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(Cancelled)
}

// NewDataLoss creates a new core error with code DataLoss and message formatted
// according to the provided format specifier.
func NewDataLoss(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(DataLoss)
}

// NewDeadlineExceeded creates a new core error with code DeadlineExceeded and message
// formatted according to the provided format specifier.
func NewDeadlineExceeded(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(DeadlineExceeded)
}

// NewFailedPrecondition creates a new core error with code FailedPrecondition and message
// formatted according to the provided format specifier.
func NewFailedPrecondition(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(FailedPrecondition)
}

// NewInternal creates a new core error with code Internal and message formatted
// according to the provided format specifier.
func NewInternal(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(Internal)
}

// NewInvalidArgument creates a new core error with code InvalidArgument and message
// formatted according to the provided format specifier.
func NewInvalidArgument(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(InvalidArgument)
}

// NewNotFound creates a new core error with code NotFound and message formatted
// according to the provided format specifier.
func NewNotFound(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(NotFound)
}

// NewPermissionDenied creates a new core error with code PermissionDenied and message
// formatted according to the provided format specifier.
func NewPermissionDenied(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(PermissionDenied)
}

// NewUnauthenticated creates a new core error with code Unauthenticated and message
// formatted according to the provided format specifier.
func NewUnauthenticated(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(Unauthenticated)
}

// NewUnavailable creates a new core error with code Unavailable and message
// formatted according to the provided format specifier.
func NewUnavailable(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(Unavailable)
}

// NewUnknown creates a new core error with code Unknown and message formatted
// according to the provided format specifier.
func NewUnknown(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(Unknown)
}

// NewUnimplemented creates a new core error with code Unimplemented and message
// formatted according to the provided format specifier.
func NewUnimplemented(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(Unimplemented)
}

// NewOutOfRange creates a new core error with code OutOfRange and message
// formatted according to the provided format specifier.
func NewOutOfRange(format string, args ...any) *Error {
	return Wrapf(format, args...).WithCode(OutOfRange)
}

// Is determines if error has code.
func Is(code Code, err error) bool {
	if code == OK && err == nil {
		return true
	}

	switch err := err.(type) {
	case *Error:
		return err.Code == code
	default:
		if st, ok := status.FromError(err); ok {
			return Code(st.Code()) == code
		}
		return false
	}
}

// IsInvalidArgument determines if err has code InvalidArgument.
func IsInvalidArgument(err error) bool {
	return Is(InvalidArgument, err)
}

// IsNotFound determines if err has code NotFound.
func IsNotFound(err error) bool {
	return Is(NotFound, err)
}

// IsAlreadyExists determines if err has code AlreadyExists.
func IsAlreadyExists(err error) bool {
	return Is(AlreadyExists, err)
}

// Map represents the context or an error.
type Map map[string]any

// Error represents an error with a code, message, and context.
type Error struct {
	Code    Code   `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Context Map    `json:"context,omitempty"`
}

// WithCode adds given code to Error.
func (m *Error) WithCode(code Code) *Error {
	m.Code = code

	return m
}

// WithErrCode adds given code from the provided error to Error.
func (m *Error) WithErrCode(err error) *Error {
	m.Code = Unknown

	if err, ok := err.(*Error); ok {
		m.Code = err.Code
	}

	return m
}

// Error returns error message or "<nil>" if none.
func (m *Error) Error() string {
	if m != nil {
		return m.Msg
	}

	return "<nil>"
}

// GetCode returns the error code.
func (m *Error) GetCode() Code {
	if m != nil {
		return m.Code
	}

	return 0
}

// GetMsg returns the error message.
func (m *Error) GetMsg() string {
	if m != nil {
		return m.Msg
	}

	return ""
}
