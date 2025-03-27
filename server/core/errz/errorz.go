package errz

import (
	"github.com/jcfug8/daylear/server/core/mapz"
)

// Context returns an errz with context.
func Context(namespace string) *Errorz {
	return &Errorz{
		Context:   make(Map),
		Namespace: namespace,
	}
}

// Errorz represents an error with a context and a namespace.
type Errorz struct {
	Context   mapz.Map
	Namespace string
}

func (errz *Errorz) nsKey(key string) string {
	if key != "" {
		return errz.Namespace + "." + key
	}

	return errz.Namespace
}

// SetContext sets an errz context value.
func (errz *Errorz) SetContext(key string, value any) *Errorz {
	mapz.Set(errz.Context, errz.nsKey(key), value)

	return errz
}

// GetContext gets an errz context value.
func (errz *Errorz) GetContext(key string) any {
	return mapz.Get(errz.Context, errz.nsKey(key))
}

// WithContext copies a error context into the errz context.
func (errz *Errorz) WithContext(err *Error) *Error {
	if err == nil {
		return nil
	}

	if err.Context == nil {
		err.Context = make(mapz.Map)
	}

	for k, v := range mapz.CopyMap(errz.Context) {
		err.Context[k] = v
	}

	return err
}

// Wrap returns a core error with code Unknown.
func (errz *Errorz) Wrap(err error) *Error {
	return errz.WithContext(Wrap(err))
}

// Wrapf returns a core error with code Wrapf and message formatted
// according to the provided format specifier.
func (errz *Errorz) Wrapf(format string, args ...any) *Error {
	return errz.WithContext(Wrapf(format, args...))
}

// NewInternal returns a core error with code Internal and message formatted
// according to the provided format specifier.
func (errz *Errorz) NewInternal(format string, args ...any) *Error {
	return errz.WithContext(NewInternal(format, args...))
}

// NewUnknown returns a core error with code Unknown and message formatted
// according to the provided format specifier.
func (errz *Errorz) NewUnknown(format string, args ...any) *Error {
	return errz.WithContext(NewUnknown(format, args...))
}

// NewNotFound returns a core error with code NotFound and message formatted
// according to the provided format specifier.
func (errz *Errorz) NewNotFound(format string, args ...any) *Error {
	return errz.WithContext(NewNotFound(format, args...))
}

// NewAlreadyExists returns a core error with code AlreadyExists and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewAlreadyExists(format string, args ...any) *Error {
	return errz.WithContext(NewAlreadyExists(format, args...))
}

// NewInvalidArgument returns a core error with code InvalidArgument and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewInvalidArgument(format string, args ...any) *Error {
	return errz.WithContext(NewInvalidArgument(format, args...))
}

// NewFailedPrecondition returns an errz with context FailedPrecondition.
func (errz *Errorz) NewFailedPrecondition(format string, args ...any) *Error {
	return errz.WithContext(NewFailedPrecondition(format, args...))
}

// NewCancelled return an errz with context Cancelled.
func (errz *Errorz) NewCancelled(format string, args ...any) *Error {
	return errz.WithContext(NewCancelled(format, args...))
}

// NewDataLoss returns an errz with context DataLoss.
func (errz *Errorz) NewDataLoss(format string, args ...any) *Error {
	return errz.WithContext(NewDataLoss(format, args...))
}

// NewDeadlineExceeded returns an errz with context DeadlineExceeded.
func (errz *Errorz) NewDeadlineExceeded(format string, args ...any) *Error {
	return errz.WithContext(NewDeadlineExceeded(format, args...))
}

// NewPermissionDenied returns a core error with code PermissionDenied and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewPermissionDenied(format string, args ...any) *Error {
	return errz.WithContext(NewPermissionDenied(format, args...))
}

// NewUnauthenticated returns a core error with code Unauthenticated and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewUnauthenticated(format string, args ...any) *Error {
	return errz.WithContext(NewUnauthenticated(format, args...))
}

// NewUnavailable returns a core error with code Unavailable and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewUnavailable(format string, args ...any) *Error {
	return errz.WithContext(NewUnavailable(format, args...))
}

// NewUnimplemented returns a core error with code Unimplemented and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewUnimplemented(format string, args ...any) *Error {
	return errz.WithContext(NewUnimplemented(format, args...))
}

// NewOutOfRange returns a core error with code OutOfRange and message
// formatted according to the provided format specifier.
func (errz *Errorz) NewOutOfRange(format string, args ...any) *Error {
	return errz.WithContext(NewOutOfRange(format, args...))
}

// IsOK returns true when error has code OK.
func (*Errorz) IsOK(err error) bool {
	return Is(OK, err)
}

// IsCancelled returns true when error has code Cancelled.
func (*Errorz) IsCancelled(err error) bool {
	return Is(Cancelled, err)
}

// IsUnknown returns true when error has code Unknown.
func (*Errorz) IsUnknown(err error) bool {
	return Is(Unknown, err)
}

// IsInvalidArgument returns true when error has code InvalidArgument.
func (*Errorz) IsInvalidArgument(err error) bool {
	return Is(InvalidArgument, err)
}

// IsDeadlineExceeded returns true when error has code DeadlineExceeded.
func (*Errorz) IsDeadlineExceeded(err error) bool {
	return Is(DeadlineExceeded, err)
}

// IsNotFound returns true when error has code NotFound.
func (*Errorz) IsNotFound(err error) bool {
	return Is(NotFound, err)
}

// IsAlreadyExists returns true when error has code AlreadyExists.
func (*Errorz) IsAlreadyExists(err error) bool {
	return Is(AlreadyExists, err)
}

// IsPermissionDenied returns true when error has code PermissionDenied.
func (*Errorz) IsPermissionDenied(err error) bool {
	return Is(PermissionDenied, err)
}

// IsResourceExhausted returns true when error has code ResourceExhausted.
func (*Errorz) IsResourceExhausted(err error) bool {
	return Is(ResourceExhausted, err)
}

// IsFailedPrecondition returns true when error has code FailedPrecondition.
func (*Errorz) IsFailedPrecondition(err error) bool {
	return Is(FailedPrecondition, err)
}

// IsAborted returns true when error has code Aborted.
func (*Errorz) IsAborted(err error) bool {
	return Is(Aborted, err)
}

// IsOutOfRange returns true when error has code OutOfRange.
func (*Errorz) IsOutOfRange(err error) bool {
	return Is(OutOfRange, err)
}

// IsUnimplemented returns true when error has code Unimplemented.
func (*Errorz) IsUnimplemented(err error) bool {
	return Is(Unimplemented, err)
}

// IsInternal returns true when error has code Internal.
func (*Errorz) IsInternal(err error) bool {
	return Is(Internal, err)
}

// IsUnavailable returns true when error has code Unavailable.
func (*Errorz) IsUnavailable(err error) bool {
	return Is(Unavailable, err)
}

// IsDataLoss returns true when error has code DataLoss.
func (*Errorz) IsDataLoss(err error) bool {
	return Is(DataLoss, err)
}

// IsUnauthenticated returns true when error has code Unauthenticated.
func (*Errorz) IsUnauthenticated(err error) bool {
	return Is(Unauthenticated, err)
}
