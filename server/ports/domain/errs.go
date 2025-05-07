package domain

// ErrNotFound is an error that indicates that a resource was not found.
type ErrNotFound struct { Msg string }

// Error returns the error message.
func (e ErrNotFound) Error() string { return e.Msg }

// ErrInternal is an error that indicates that an internal error occurred.
type ErrInternal struct { Msg string }

// Error returns the error message.
func (e ErrInternal) Error() string { return e.Msg }
