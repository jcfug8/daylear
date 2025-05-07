package domain

// ErrNotFound is an error that indicates that a resource was not found.
type ErrNotFound struct{ Msg string }

// Error returns the error message.
func (e ErrNotFound) Error() string { return e.Msg }

// ErrInternal is an error that indicates that an internal error occurred.
type ErrInternal struct{ Msg string }

// Error returns the error message.
func (e ErrInternal) Error() string { return e.Msg }

// ErrInvalidArgument is an error that indicates that an invalid argument was provided.
type ErrInvalidArgument struct{ Msg string }

// Error returns the error message.
func (e ErrInvalidArgument) Error() string { return e.Msg }

// ErrPermissionDenied is an error that indicates that a user does not have permission to access a resource.
type ErrPermissionDenied struct{ Msg string }

// Error returns the error message.
func (e ErrPermissionDenied) Error() string { return e.Msg }
