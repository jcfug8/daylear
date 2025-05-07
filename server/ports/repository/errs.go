package repository

// ErrNotFound is an error that indicates that a resource was not found.
type ErrNotFound struct { Msg string }

// Error returns the error message.
func (e ErrNotFound) Error() string { return e.Msg }

// ErrInternal is an error that indicates that an internal error occurred.
type ErrInternal struct { Msg string }

// Error returns the error message.
func (e ErrInternal) Error() string { return e.Msg }

// ErrInvalidArgument is an error that indicates that an invalid argument was provided.
type ErrInvalidArgument struct { Msg string }

// Error returns the error message.
func (e ErrInvalidArgument) Error() string { return e.Msg }

// ErrNewAlreadyExists is an error that indicates that a resource already exists.
type ErrNewAlreadyExists struct { Msg string }

// Error returns the error message.
func (e ErrNewAlreadyExists) Error() string { return e.Msg }

// ErrNewUnimplemented is an error that indicates that a resource is not implemented.
type ErrNewUnimplemented struct { Msg string }

// Error returns the error message.
func (e ErrNewUnimplemented) Error() string { return e.Msg }

// ErrNewNotFound is an error that indicates that a resource was not found.
type ErrNewNotFound struct { Msg string }

// Error returns the error message.
func (e ErrNewNotFound) Error() string { return e.Msg }
